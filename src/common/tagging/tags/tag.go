// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package tags

import (
	"fmt"
	"regexp"
)

type Tag struct {
	Key   string
	Value string
}

const TraceTagKey = "dd_correlation_uuid"
const GitFileTagKey = "dd_git_file"
const GitModifiersTagKey = "dd_git_modifiers"
const GitLastModifiedAtTagKey = "dd_git_last_modified_at"
const GitLastModifiedByTagKey = "dd_git_last_modified_by"
const GitCreatedAtTagKey = "dd_git_created_at"
const GitCreatedByTagKey = "dd_git_created_by"
const GitRepoTagKey = "dd_git_repo"
const GitRepoUrlTagKey = "dd_git_repo_url"
const ResourceNameTagKey = "dd_resource_name"
const ResourceSignatureTagKey = "dd_git_resource_signature"

type ITag interface {
	Init()
	CalculateValue(data interface{}) (ITag, error)
	GetKey() string
	SetValue(val string)
	GetValue() string
	GetPriority() int
	GetDescription() string
	SetTagPrefix(tagPrefix string)
}

type TagDiff struct {
	Key       string
	PrevValue string
	NewValue  string
}

func Init(key string, value string) ITag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func (t *Tag) Init() {}

func (t *Tag) SetTagPrefix(tagPrefix string) {
	t.Key = fmt.Sprintf("%s%s", tagPrefix, t.Key)
}

func (t *Tag) GetPriority() int {
	return 0
}

func (t *Tag) CalculateValue(_ interface{}) (ITag, error) {
	return &Tag{
		Key:   t.Key,
		Value: t.Value,
	}, nil
}

func (t *Tag) GetDescription() string {
	return "Abstract tag class"
}

func (t *Tag) GetKey() string {
	return t.Key
}

func (t *Tag) SetValue(val string) {
	t.Value = val
}

func (t *Tag) GetValue() string {
	return t.Value
}

// IsTagKeyMatch Try to match the tag's key name with a potentially quoted string
func IsTagKeyMatch(tag ITag, keyName string) bool {
	match, _ := regexp.Match(fmt.Sprintf(`\b"?%s"?\b`, regexp.QuoteMeta(keyName)), []byte(tag.GetKey()))
	return match
}
