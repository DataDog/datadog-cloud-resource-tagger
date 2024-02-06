// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/Datadog/cloud-resource-tagger/src/common"
	"github.com/Datadog/cloud-resource-tagger/src/common/clioptions"
	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	"github.com/Datadog/cloud-resource-tagger/src/common/reports"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging"
	taggingUtils "github.com/Datadog/cloud-resource-tagger/src/common/tagging/utils"
	"github.com/Datadog/cloud-resource-tagger/src/common/utils"
	tfStructure "github.com/Datadog/cloud-resource-tagger/src/terraform/structure"
)

type Runner struct {
	TagGroups            []tagging.ITagGroup
	parsers              []common.IParser
	ChangeAccumulator    *reports.TagChangeAccumulator
	reportingService     *reports.ReportService
	dir                  string
	skipDirs             []string
	skippedTags          []string
	configFilePath       string
	skippedResourceTypes []string
	skippedResources     []string
	workersNum           int
	dryRun               bool
	localModuleTag       bool
	changedFiles         []string
}

func (r *Runner) Init(commands *clioptions.TagOptions) error {
	dir := commands.Directory
	if dir == "" {
		dir = utils.DetermineTopLevelDirectory(commands.ChangedFiles)
	}
	r.parsers = append(r.parsers, &tfStructure.TerraformParser{})

	options := map[string]string{}
	for _, parser := range r.parsers {
		parser.Init(dir, options)
	}

	for _, group := range commands.TagGroups {
		tagGroup := taggingUtils.TagGroupsByName(taggingUtils.TagGroupName(group))
		r.TagGroups = append(r.TagGroups, tagGroup)
	}

	for _, tagGroup := range r.TagGroups {
		tagGroup.InitTagGroup(dir, commands.SkipTags, commands.Tag)
	}

	r.changedFiles = utils.SplitStringByComma(commands.ChangedFiles)

	r.ChangeAccumulator = reports.TagChangeAccumulatorInstance
	r.reportingService = reports.ReportServiceInst
	r.dir = commands.Directory
	r.dryRun = commands.DryRun

	r.workersNum = 10
	if utils.InSlice(r.skipDirs, r.dir) {
		logger.Warning(fmt.Sprintf("Selected dir, %s, is skipped - expect an empty result", r.dir))
	}
	return nil
}

func (r *Runner) worker(fileChan chan string, wg *sync.WaitGroup) {
	//("Starting worker")
	for file := range fileChan {
		r.TagFile(file)
		wg.Done()
	}
}

func (r *Runner) TagChangedFiles() (*reports.ReportService, error) {
	var changedFiles = r.changedFiles

	// Set up the wait group and the file channel
	var wg sync.WaitGroup
	wg.Add(len(changedFiles))
	fileChan := make(chan string)

	// Tag the changed files
	for i := 0; i < r.workersNum; i++ {
		go r.worker(fileChan, &wg)
	}

	for _, file := range changedFiles {
		fileChan <- file
	}
	close(fileChan)
	wg.Wait()

	for _, parser := range r.parsers {
		parser.Close()
	}

	return r.reportingService, nil
}

func (r *Runner) TagDirectory() (*reports.ReportService, error) {
	//fmt.Println("Starting tagging directory")

	var files []string
	err := filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to scan dir %s", path))
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logger.Error("Failed to run Walk() on root dir", r.dir)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	fileChan := make(chan string)

	for i := 0; i < r.workersNum; i++ {
		go r.worker(fileChan, &wg)
	}

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)
	wg.Wait()

	for _, parser := range r.parsers {
		parser.Close()
	}

	return r.reportingService, nil
}

func (r *Runner) isSkippedResourceType(resourceType string) bool {
	for _, skippedResourceType := range r.skippedResourceTypes {
		if resourceType == skippedResourceType {
			return true
		}
	}
	return false
}

func (r *Runner) isSkippedResource(resource string) bool {
	for _, skippedResource := range r.skippedResources {
		if resource == skippedResource {
			return true
		}
	}
	return false
}

func (r *Runner) TagFile(file string) {
	for _, parser := range r.parsers {
		if r.isFileSkipped(parser, file) {
			logger.Debug(fmt.Sprintf("%v parser Skipping %v", parser.Name(), file))
			continue
		}
		logger.Info(fmt.Sprintf("Tagging %v\n", file))

		blocks, err := parser.ParseFile(file)
		if err != nil {
			logger.Info(fmt.Sprintf("Failed to parse file %v with parser %v", file, reflect.TypeOf(parser)))
			continue
		}
		isFileTaggable := false
		for _, block := range blocks {
			if r.isSkippedResourceType(block.GetResourceType()) {
				continue
			}
			if r.isSkippedResource(block.GetResourceID()) {
				continue
			}
			if block.IsBlockTaggable() {
				logger.Debug(fmt.Sprintf("Tagging %v:%v", file, block.GetResourceID()))
				isFileTaggable = true
				for _, tagGroup := range r.TagGroups {
					err := tagGroup.CreateTagsForBlock(block)
					if err != nil {
						logger.Warning(fmt.Sprintf("Failed to tag %v in %v due to %v", block.GetResourceID(), block.GetFilePath(), err.Error()))
						continue
					}
				}
			} else {
				logger.Debug(fmt.Sprintf("Block %v:%v is not taggable, skipping", file, block.GetResourceID()))
			}
			r.ChangeAccumulator.AccumulateChanges(block)
		}
		if isFileTaggable && !r.dryRun {
			err = parser.WriteFile(file, blocks, file)
			if err != nil {
				logger.Warning(fmt.Sprintf("Failed writing tags to file %s, because %v", file, err))
			}
		}
	}
}

func (r *Runner) isFileSkipped(p common.IParser, file string) bool {
	relPath, _ := filepath.Rel(r.dir, file)
	for _, sp := range r.skipDirs {
		if strings.HasPrefix(filepath.Join(r.dir, relPath), sp) {
			return true
		}
	}

	matchingSuffix := false
	for _, suffix := range p.GetSupportedFileExtensions() {
		if strings.HasSuffix(file, suffix) {
			matchingSuffix = true
		}
	}
	if !matchingSuffix {
		return true
	}
	for _, pattern := range p.GetSkippedDirs() {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return !p.ValidFile(file)
}
