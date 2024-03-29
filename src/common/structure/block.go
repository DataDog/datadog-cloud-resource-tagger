// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package structure

import (
	"sort"

	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type Lines struct {
	Start int
	End   int
}

type TagDiff struct {
	Added   []tags.ITag
	Updated []*tags.TagDiff
}

var SpecialResourceTypes = map[string]int{
	"aws_db_proxy":      10,
	"AWS::RDS::DBProxy": 10,
}

type IBlock interface {
	Init(filePath string, rawBlock interface{})
	GetFilePath() string
	GetLines(...bool) Lines
	GetExistingTags() []tags.ITag
	GetNewTags() []tags.ITag
	GetRawBlock() interface{}
	GetTraceID() string
	AddNewTags(newTags []tags.ITag)
	MergeTags() []tags.ITag
	CalculateTagsDiff() *TagDiff
	IsBlockTaggable() bool
	GetResourceID() string
	GetTagsLines() Lines
	GetSeparator() string
	GetTagsAttributeName() string
	IsGCPBlock() bool
	GetResourceType() string
	GetResourceName() string
}

type Block struct {
	FilePath          string
	ExistingTags      []tags.ITag
	NewTags           []tags.ITag
	RawBlock          interface{}
	IsTaggable        bool
	TagsAttributeName string
	Lines             Lines
	TagLines          Lines
	Name              string
	Type              string
}

func (b *Block) Init(filePath string, rawBlock interface{}) {
	b.FilePath = filePath
	b.RawBlock = rawBlock
}

func (b *Block) GetLines(_ ...bool) Lines {
	return b.Lines
}

func (b *Block) GetResourceID() string {
	return b.Name
}

func (b *Block) GetResourceType() string {
	return b.Type
}

func (b *Block) GetTagsLines() Lines {
	return b.TagLines
}

func (b *Block) GetSeparator() string {
	panic("implement me")
}

func (b *Block) AddNewTags(newTags []tags.ITag) {
	if newTags == nil {
		return
	}
	isTraced := false
	traceTagKey := tags.TraceTagKey
	for _, tag := range b.ExistingTags {
		match := tags.IsTagKeyMatch(tag, traceTagKey)
		if match {
			isTraced = true
			break
		}
	}
	if isTraced {
		traceTagTraceIndex := -1
		for index, tag := range newTags {
			match := tags.IsTagKeyMatch(tag, traceTagKey)
			if match {
				traceTagTraceIndex = index
			}
		}

		if traceTagTraceIndex >= 0 {
			newTags = append(newTags[:traceTagTraceIndex], newTags[traceTagTraceIndex+1:]...)
		}
	}
	b.NewTags = append(b.NewTags, newTags...)
	sort.Slice(b.NewTags, func(i, j int) bool {
		return b.NewTags[i].GetKey() > b.NewTags[j].GetKey()
	})
	if limit, ok := SpecialResourceTypes[b.GetResourceType()]; ok && len(b.NewTags)+len(b.ExistingTags) > limit {
		b.NewTags = b.NewTags[0 : limit-len(b.ExistingTags)]
	}
}

// MergeTags merges the tags and returns all the tags.
func (b *Block) MergeTags() []tags.ITag {
	existingTagsByKey := map[string]tags.ITag{}
	newTagsByKey := map[string]tags.ITag{}

	for _, tag := range b.ExistingTags {
		existingTagsByKey[tag.GetKey()] = tag
	}
	for _, tag := range b.NewTags {
		newTagsByKey[tag.GetKey()] = tag
	}

	var mergedTags []tags.ITag
	traceTagKeyName := tags.TraceTagKey
	for _, existingTag := range b.ExistingTags {
		if newTag, ok := newTagsByKey[existingTag.GetKey()]; ok {
			match := tags.IsTagKeyMatch(existingTag, traceTagKeyName)
			if match {
				mergedTags = append(mergedTags, existingTag)
			} else {
				mergedTags = append(mergedTags, newTag)
			}
			delete(newTagsByKey, existingTag.GetKey())
		} else {
			mergedTags = append(mergedTags, existingTag)
		}
	}

	for newTagKey := range newTagsByKey {
		mergedTags = append(mergedTags, newTagsByKey[newTagKey])
	}

	return mergedTags
}

// CalculateTagsDiff returns a map which explains the changes in tags for this block
// Added is the new tags, Updated is the tags which were modified
func (b *Block) CalculateTagsDiff() *TagDiff {
	var diff = TagDiff{}
	for _, newTag := range b.GetNewTags() {
		found := false
		for _, existingTag := range b.GetExistingTags() {
			if newTag.GetKey() == existingTag.GetKey() {
				found = true
				if newTag.GetValue() != existingTag.GetValue() {
					diff.Updated = append(diff.Updated, &tags.TagDiff{
						Key:       newTag.GetKey(),
						PrevValue: existingTag.GetValue(),
						NewValue:  newTag.GetValue(),
					})
					break
				}
			}
		}
		if !found {
			diff.Added = append(diff.Added, newTag)
		}
	}
	return &diff
}

func (b *Block) GetRawBlock() interface{} {
	return b.RawBlock
}

func (b *Block) GetExistingTags() []tags.ITag {
	return b.ExistingTags
}

func (b *Block) GetNewTags() []tags.ITag {
	return b.NewTags
}

func (b *Block) IsBlockTaggable() bool {
	return b.IsTaggable
}

func (b *Block) GetFilePath() string {
	return b.FilePath
}

func (b *Block) GetTraceID() string {
	for _, tag := range b.MergeTags() {
		if tag.GetKey() == tags.TraceTagKey {
			return tag.GetValue()
		}
	}
	return ""
}

func (b *Block) GetTagsAttributeName() string {
	return b.TagsAttributeName
}

func (b *Block) IsGCPBlock() bool {
	return false
}

func (b *Block) GetResourceName() string {
	return b.GetResourceID()
}
