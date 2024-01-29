package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Datadog/cloud-resource-tagger/src/common"
	"github.com/Datadog/cloud-resource-tagger/src/common/clioptions"
	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	"github.com/Datadog/cloud-resource-tagger/src/common/reports"
	"github.com/Datadog/cloud-resource-tagger/src/common/runner"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/utils"
)

func main() {
	app := &cli.App{
		Name:                   "cloud-resource-tagger",
		HelpName:               "",
		Usage:                  "Tag your cloud resources with git information",
		Version:                common.Version,
		Description:            "Cloud resource tagger",
		Compiled:               time.Time{},
		Authors:                []*cli.Author{{Name: "Datadog", Email: "bahar.shah@datadoghq.com"}},
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			tagCommand(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err.Error())
	}
}

func tagCommand() *cli.Command {
	directoryArg := "directory"
	tagArg := "tags"
	outputArg := "output"
	outputJSONFileArg := "output-json-file"
	tagGroupArg := "tag-groups"
	dryRunArgs := "dry-run"

	return &cli.Command{
		Name:                   "tag",
		Usage:                  "apply datadog cloud resource tagging across your files",
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		Action: func(c *cli.Context) error {
			options := clioptions.TagOptions{
				Directory:      c.String(directoryArg),
				Tag:            c.StringSlice(tagArg),
				Output:         c.String(outputArg),
				OutputJSONFile: c.String(outputJSONFileArg),
				TagGroups:      c.StringSlice(tagGroupArg),
				DryRun:         c.Bool(dryRunArgs),
			}
			options.Validate()

			return tag(&options)
		},
		Flags: []cli.Flag{ // When adding flags, make sure they are supported in the GitHub action as well via entrypoint.sh
			&cli.StringFlag{
				Name:        directoryArg,
				Aliases:     []string{"d"},
				Usage:       "directory to tag",
				Required:    true,
				DefaultText: "path/to/iac/root",
			},
			&cli.StringFlag{
				Name:        outputArg,
				Aliases:     []string{"o"},
				Usage:       "set output format",
				Value:       "cli",
				DefaultText: "json",
			},
			&cli.StringSliceFlag{
				Name:        tagGroupArg,
				Aliases:     []string{"g"},
				Usage:       "Narrow down results to the matching tag groups",
				Value:       cli.NewStringSlice(utils.GetAllTagGroupsNames()...),
				DefaultText: "git,code2cloud",
			},
			&cli.StringFlag{
				Name:        outputJSONFileArg,
				Usage:       "json file path for output",
				DefaultText: "result.json",
			},
		},
	}
}

func tag(options *clioptions.TagOptions) error {
	iacTaggerRunner := new(runner.Runner)

	logger.Info(fmt.Sprintf("Setting up to tag the directory %v\n", options.Directory))
	err := iacTaggerRunner.Init(options)
	if err != nil {
		logger.Error(err.Error())
	}
	reportService, err := iacTaggerRunner.TagDirectory()
	if err != nil {
		logger.Error(err.Error())
	}
	printReport(reportService, options)

	return nil
}

func printReport(reportService *reports.ReportService, options *clioptions.TagOptions) {
	reportService.CreateReport()

	if options.OutputJSONFile != "" {
		reportService.PrintJSONToFile(options.OutputJSONFile)
	}

	switch strings.ToLower(options.Output) {
	case "cli":
		reportService.PrintToStdout()
	case "json":
		reportService.PrintJSONToStdout()
	default:
		return
	}
}
