// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package clioptions

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/validator.v2"

	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	taggingUtils "github.com/Datadog/cloud-resource-tagger/src/common/tagging/utils"
	"github.com/Datadog/cloud-resource-tagger/src/common/utils"
)

var allowedOutputTypes = []string{"cli", "json"}

type TagOptions struct {
	Directory            string
	Output               string `validate:"output"`
	OutputJSONFile       string
	Tag                  []string
	SkipTags             []string
	TagGroups            []string `validate:"tagGroupNames"`
	DryRun               bool
	ChangedFiles         []string
	IncludeResourceTypes []string
}

type ListTagsOptions struct {
	TagGroups []string `validate:"tagGroupNames"`
}

func (o *TagOptions) Validate() {
	_ = validator.SetValidationFunc("output", validateOutput)
	_ = validator.SetValidationFunc("tagGroupNames", validateTagGroupNames)

	o.Tag = utils.SplitStringByComma(o.Tag)
	o.TagGroups = utils.SplitStringByComma(o.TagGroups)
	o.IncludeResourceTypes = utils.SplitStringByComma(o.IncludeResourceTypes)

	if err := validator.Validate(o); err != nil {
		logger.Error(err.Error())
	}
}

func (l *ListTagsOptions) Validate() {
	_ = validator.SetValidationFunc("tagGroupNames", validateTagGroupNames)
	l.TagGroups = utils.SplitStringByComma(l.TagGroups)

	if err := validator.Validate(l); err != nil {
		logger.Error(err.Error())
	}
}

func validateTagGroupNames(v interface{}, _ string) error {
	tagGroupsNames := taggingUtils.GetAllTagGroupsNames()
	val, ok := v.([]string)
	if ok {
		for _, gn := range val {
			if !utils.InSlice(tagGroupsNames, gn) {
				return fmt.Errorf("tag group %s is not one of the supported tag groups. supported groups: %v", gn, tagGroupsNames)
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported tag group names [%s]. supported types: %s", val, tagGroupsNames)
}

func validateOutput(v interface{}, _ string) error {
	val, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}

	if val != "" && !utils.InSlice(allowedOutputTypes, strings.ToLower(val)) {
		return fmt.Errorf("unsupported output type [%s]. allowed types: %s", val, allowedOutputTypes)
	}

	return nil
}

func validateConfigFile(v interface{}, _ string) error {
	if v != "" {
		val, ok := v.(string)
		if !ok {
			return validator.ErrUnsupported
		}

		if _, err := os.Stat(val); err != nil {
			return fmt.Errorf("configuration file %s does not exist", v)
		}

	}
	return nil
}
