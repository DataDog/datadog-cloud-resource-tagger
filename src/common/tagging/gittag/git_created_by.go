// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gittag

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitCreatedByTag struct {
	tags.Tag
}

func (t *GitCreatedByTag) Init() {
	t.Key = tags.GitCreatedByTagKey
}

func (t *GitCreatedByTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	firstCommit := gitBlame.GetFirstCommit()

	if firstCommit == nil {
		return nil, fmt.Errorf("first commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: firstCommit.Author}, nil
}

func (t *GitCreatedByTag) GetDescription() string {
	return "Who created this resource's configuration"
}
