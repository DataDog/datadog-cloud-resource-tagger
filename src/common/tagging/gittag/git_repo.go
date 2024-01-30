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

type GitRepoTag struct {
	tags.Tag
}

func (t *GitRepoTag) Init() {
	t.Key = tags.GitRepoTagKey
}

func (t *GitRepoTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &tags.Tag{Key: t.Key, Value: gitBlame.GitRepository}, nil
}

func (t *GitRepoTag) GetDescription() string {
	return "The repository where this resource is provisioned in IaC"
}
