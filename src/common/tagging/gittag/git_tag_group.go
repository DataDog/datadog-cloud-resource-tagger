// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gittag

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	"github.com/Datadog/cloud-resource-tagger/src/common/structure"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
	"github.com/Datadog/cloud-resource-tagger/src/common/utils"
)

type TagGroup struct {
	tagging.TagGroup
	GitService *gitservice.GitService
}

type fileLineMapper struct {
	originToGit map[int]int
	gitToOrigin map[int]int
}

func (t *TagGroup) InitTagGroup(path string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	opt := tagging.InitTagGroupOptions{}
	for _, fn := range options {
		fn(&opt)
	}
	t.SkippedTags = skippedTags
	t.SpecifiedTags = explicitlySpecifiedTags
	t.Options = opt
	gitService, err := gitservice.NewGitService(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize git service for path \"%s\". Please ensure the provided root directory is initialized via the git init command: %q", path, err), "SILENT")
	}
	t.GitService = gitService

	t.SetTags(t.GetDefaultTags())
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&GitOrgTag{},
		&GitRepoTag{},
		&GitRepoUrlTag{},
		&GitFileTag{},
		&GitModifiedCommitTag{},
		&GitModifiersTag{},
		&GitLastModifiedAtTag{},
		&GitLastModifiedByTag{},
		&GitCreatedAtTag{},
		&GitCreatedByTag{},
		&GitCreateCommitTag{},
		&GitResourceLinesTag{},
	}
}

func (t *TagGroup) initFileMapping(path string) fileLineMapper {
	fileBlame, err := t.GitService.GetFileBlame(path)
	if err != nil {
		logger.Warning(fmt.Sprintf("Unable to get git blame for file %s: %s", path, err))
		return fileLineMapper{}
	}

	return t.mapOriginFileToGitFile(path, fileBlame)
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	fileLinesMap := t.initFileMapping(block.GetFilePath())
	linesInGit := t.getBlockLinesInGit(block, fileLinesMap)

	blame, err := t.GitService.GetBlameForFileLines(block.GetFilePath(), linesInGit)

	if err != nil || blame == nil {
		logger.Warning(fmt.Sprintf("Failed to tag %v with git tags, err: %v", block.GetResourceID(), err.Error()))

		org := t.GitService.GetOrganization()
		repo := t.GitService.GetRepoName()
		blame = &gitservice.GitBlame{
			GitOrg:           org,
			GitRepository:    repo,
			GitRepositoryUrl: fmt.Sprintf("git@github.com:%s/%s", org, repo),
			FilePath:         block.GetFilePath(),
			AbsoluteFilePath: block.GetFilePath(),
		}
		logger.Warning(fmt.Sprintf("Created manual git blame for %s", block.GetFilePath()))
	}

	// if !t.hasNonTagChanges(blame, block) {
	// 	return nil
	// }

	err = t.UpdateBlockTags(block, blame)
	if err != nil {
		return err
	}
	if block.IsGCPBlock() {
		for _, tag := range block.GetNewTags() {
			t.cleanGCPTagValue(tag)
		}
	}
	return nil
}

func (t *TagGroup) getBlockLinesInGit(block structure.IBlock, linesMap fileLineMapper) structure.Lines {
	blockLines := block.GetLines()
	originToGit := linesMap.originToGit
	originStart := blockLines.Start
	originEnd := blockLines.End
	gitStart := -1
	gitEnd := -1

	for gitStart == -1 && originStart <= originEnd {
		// find the first mapped line
		gitStart = originToGit[originStart]
		originStart++
	}

	for gitEnd == -1 && originEnd >= blockLines.Start {
		// find the last mapped line
		gitEnd = originToGit[originEnd]
		originEnd--
	}

	return structure.Lines{Start: gitStart, End: gitEnd}
}

// The function maps between the scanned file lines to the lines in the git blame
func (t *TagGroup) mapOriginFileToGitFile(path string, fileBlame *git.BlameResult) fileLineMapper {
	mapper := fileLineMapper{
		originToGit: make(map[int]int),
		gitToOrigin: make(map[int]int),
	}

	gitLines := make([]string, 0)
	for _, line := range fileBlame.Lines {
		gitLines = append(gitLines, line.Text)
	}

	originFileText, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fileLineMapper{}
	}

	originLines := utils.GetLinesFromBytes(originFileText)

	matcher := difflib.NewMatcher(originLines, gitLines)
	matches := matcher.GetMatchingBlocks()
	currOriginStart := 0
	currGitStart := 0
	for _, match := range matches {
		startInOrigin := match.A
		startInGit := match.B

		// if there were lines that weren't in the range of the previous block, they are changed in the opposite file and will set to -1
		for i := currOriginStart + 1; i <= startInOrigin; i++ {
			mapper.originToGit[i] = -1
		}
		for i := currGitStart + 1; i <= startInGit; i++ {
			mapper.gitToOrigin[i] = -1
		}

		// iterate the matching block and map the corresponding lines
		for i := 1; i <= match.Size; i++ {
			mapper.originToGit[startInOrigin+i] = startInGit + i
			mapper.gitToOrigin[startInGit+i] = startInOrigin + i
		}

		currOriginStart = startInOrigin + match.Size
		currGitStart = startInGit + match.Size
	}

	return mapper
}

func (t *TagGroup) updateBlameForOriginLines(block structure.IBlock, blame *gitservice.GitBlame, fileMapping map[int]int) {
	gitBlameLines := blame.BlamesByLine
	blockLines := block.GetLines(true)
	newBlameByLines := make(map[int]*git.Line)

	for blockLine := blockLines.Start; blockLine <= blockLines.End; blockLine++ {
		if fileMapping[blockLine] == -1 {
			newBlameByLines[blockLine] = &git.Line{
				Author: blame.GitUserEmail,
				Date:   time.Now().UTC(),
				Hash:   plumbing.ZeroHash,
			}
		} else {
			newBlameByLines[blockLine] = gitBlameLines[fileMapping[blockLine]]
		}
	}

	blame.BlamesByLine = newBlameByLines
}

func (t *TagGroup) hasNonTagChanges(blame *gitservice.GitBlame, block structure.IBlock) bool {
	tagsLines := block.GetTagsLines()
	hasTags := tagsLines.Start != -1 && tagsLines.End != -1

	// if have no tags, all tags need to be added
	if !hasTags {
		return true
	}
	for lineNum, line := range blame.BlamesByLine {
		// not in the scope of tags lines, no need to compare
		if lineNum >= tagsLines.Start && lineNum <= tagsLines.End {
			if line.Hash.String() == blame.GetLatestCommitForLine(lineNum).Hash.String() {
				return true
			}
		}
	}
	return false
}

func (t *TagGroup) cleanGCPTagValue(val tags.ITag) {
	updated := val.GetValue()
	switch val.GetKey() {
	case tags.GitModifiersTagKey:
		modifiers := strings.Split(updated, "/")
		for i, m := range modifiers {
			modifiers[i] = utils.RemoveGcpInvalidChars.ReplaceAllString(m, "")
		}
		updated = strings.Join(modifiers, "__")
	case tags.GitLastModifiedAtTagKey:
		updated = strings.ReplaceAll(updated, " ", "-")
		updated = strings.ReplaceAll(updated, ":", "-")
	case tags.GitFileTagKey:
		updated = strings.ReplaceAll(updated, "/", "__")
		updated = strings.ReplaceAll(updated, ".", "_")
	case tags.GitLastModifiedByTagKey:
		updated = strings.Split(updated, "@")[0]
		updated = utils.RemoveGcpInvalidChars.ReplaceAllString(updated, "")
	case tags.GitRepoTagKey:
		updated = strings.ReplaceAll(updated, "/", "__")
		updated = strings.ReplaceAll(updated, ".", "_")
	case tags.GitRepoUrlTagKey:
		updated = strings.ReplaceAll(updated, "/", "__")
		updated = strings.ReplaceAll(updated, ".", "_")
	}

	val.SetValue(updated)
}

func GetMinimumDefaultGitTags() []string {
	return []string{"dd_git_org,dd_git_repo,dd_git_file,dd_git_resource_signature"}
}
