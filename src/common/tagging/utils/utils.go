// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package utils

import (
	"sort"

	"github.com/Datadog/cloud-resource-tagger/src/common/tagging"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/code2cloud"

	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/gittag"
)

type TagGroupName string

const (
	SimpleTagGroupName TagGroupName = "simple"
	GitTagGroupName    TagGroupName = "git"
	Code2Cloud         TagGroupName = "code2cloud"
	//ExternalTagName    TagGroupName = "external"
)

var tagGroupsByName = map[TagGroupName]tagging.ITagGroup{
	//SimpleTagGroupName: &simple.TagGroup{},
	GitTagGroupName: &gittag.TagGroup{},
	Code2Cloud:      &code2cloud.TagGroup{},
}

func TagGroupsByName(name TagGroupName) tagging.ITagGroup {
	var tagGroup tagging.ITagGroup
	switch name {
	//case SimpleTagGroupName:
	//	tagGroup = &simple.TagGroup{}
	case GitTagGroupName:
		tagGroup = &gittag.TagGroup{}
	case Code2Cloud:
		tagGroup = &code2cloud.TagGroup{}
		//case ExternalTagName:
		//	tagGroup = &external.TagGroup{}
	}

	return tagGroup
}

func GetAllTagGroupsNames() []string {
	tagGroupNames := make([]string, 0)

	for name := range tagGroupsByName {
		tagGroupNames = append(tagGroupNames, string(name))
	}
	sort.Strings(tagGroupNames)
	//tagGroupNames = append(tagGroupNames, string(ExternalTagName)) // Add the external tag name as the last tag group
	return tagGroupNames
}
