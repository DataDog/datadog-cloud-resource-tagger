// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gittag

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitResourceLinesTag struct {
	tags.Tag
}

func (t *GitResourceLinesTag) Init() {
	t.Key = "dd_git_resource_lines"
}

func (t *GitResourceLinesTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	var keys []int
	for key := range gitBlame.BlamesByLine {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	var firstLine = keys[0]
	var lastLine = -1
	if len(keys) > 0 {
		lastLine = keys[len(keys)-1]
	}
	if gitBlame.BlamesByLine == nil {
		return &tags.Tag{Key: t.Key, Value: ""}, fmt.Errorf("Failed to get resource lines")
	}
	return &tags.Tag{Key: t.Key, Value: strconv.Itoa(firstLine) + ":" + strconv.Itoa(lastLine)}, nil
}

func (t *GitResourceLinesTag) GetDescription() string {
	return "Range of lines where the resource is defined in the file"
}
