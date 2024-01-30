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

type GitLastModifiedAtTag struct {
	tags.Tag
}

func (t *GitLastModifiedAtTag) Init() {
	t.Key = tags.GitLastModifiedAtTagKey
}

func (t *GitLastModifiedAtTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	latestCommit := gitBlame.GetLatestCommit()

	if latestCommit == nil {
		return nil, fmt.Errorf("latest commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: latestCommit.Date.UTC().Format("2006-01-02 15:04:05")}, nil
}

func (t *GitLastModifiedAtTag) GetDescription() string {
	return "The last time this resource's configuration was modified"
}
