// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gittag

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitModifiersTag struct {
	tags.Tag
}

func (t *GitModifiersTag) Init() {
	t.Key = tags.GitModifiersTagKey
}

func (t *GitModifiersTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	foundModifyingUsers := make(map[string]bool)
	var modifyingUsers []string
	for _, v := range gitBlame.BlamesByLine {
		userName := strings.Split(v.Author, "@")[0]
		if !foundModifyingUsers[userName] && userName != "" && !strings.Contains(userName, "[") {
			modifyingUsers = append(modifyingUsers, userName)
			foundModifyingUsers[userName] = true
		}
	}

	sort.Strings(modifyingUsers)

	return &tags.Tag{Key: t.Key, Value: strings.Join(modifyingUsers, "/")}, nil
}

func (t *GitModifiersTag) GetDescription() string {
	return "The users who modified this resource"
}
