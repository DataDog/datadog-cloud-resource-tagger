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

type GitFileTag struct {
	tags.Tag
}

func (t *GitFileTag) Init() {
	t.Key = tags.GitFileTagKey
}

func (t *GitFileTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	fmt.Println("GitFileTag: ", gitBlame.AbsoluteFilePath)
	return &tags.Tag{Key: t.Key, Value: gitBlame.AbsoluteFilePath}, nil
}

func (t *GitFileTag) GetDescription() string {
	return "The file (including path) in the repository where this resource is provisioned in IaC"
}
