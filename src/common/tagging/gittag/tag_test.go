// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gittag

import (
	"os"
	"testing"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
	"github.com/Datadog/cloud-resource-tagger/tests/utils/blameutils"

	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	blame := blameutils.SetupBlame(t)
	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := GitOrgTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_org")
		assert.Equal(t, blameutils.Org, valueTag.GetValue())
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := GitRepoTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_repo")
		assert.Equal(t, blameutils.Repository, valueTag.GetValue())
	})

	t.Run("GitRepoUrlTagCreation", func(t *testing.T) {
		tag := GitRepoUrlTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_repo_url")
		assert.Equal(t, blameutils.RepositoryUrl, valueTag.GetValue())
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := GitFileTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_file")
		assert.Equal(t, blameutils.AbsoluteFilePath, valueTag.GetValue())
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := GitModifiedCommitTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_modified_commit")
		assert.Equal(t, blameutils.CommitHash1, valueTag.GetValue())
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := GitLastModifiedAtTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, valueTag.GetKey(), "dd_git_last_modified_at")
		assert.Equal(t, "2020-03-28 21:42:46", valueTag.GetValue())
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := GitLastModifiedByTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "dd_git_last_modified_by", valueTag.GetKey())
		assert.Equal(t, "a@a.com", valueTag.GetValue())
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := GitModifiersTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "dd_git_modifiers", valueTag.GetKey())
		assert.Equal(t, "a", valueTag.GetValue())
	})

	t.Run("Tag description tests", func(t *testing.T) {
		tag := tags.Tag{}
		defaultDescription := tag.GetDescription()
		cwd, _ := os.Getwd()
		g := TagGroup{}
		g.InitTagGroup(cwd, nil, nil)
		for _, tag := range g.GetTags() {
			assert.NotEqual(t, defaultDescription, tag.GetDescription())
			assert.NotEqual(t, "", tag.GetDescription())
		}
	})

}

func TestTagCreationWithPrefix(t *testing.T) {
	blame := blameutils.SetupBlame(t)
	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := GitOrgTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_org", valueTag.GetKey())
		assert.Equal(t, blameutils.Org, valueTag.GetValue())
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := GitRepoTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_repo", valueTag.GetKey())
		assert.Equal(t, blameutils.Repository, valueTag.GetValue())
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := GitFileTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_file", valueTag.GetKey())
		assert.Equal(t, blameutils.AbsoluteFilePath, valueTag.GetValue())
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := GitModifiedCommitTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_modified_commit", valueTag.GetKey())
		assert.Equal(t, blameutils.CommitHash1, valueTag.GetValue())
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := GitLastModifiedAtTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_last_modified_at", valueTag.GetKey())
		assert.Equal(t, "2020-03-28 21:42:46", valueTag.GetValue())
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := GitLastModifiedByTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_last_modified_by", valueTag.GetKey())
		assert.Equal(t, "a@a.com", valueTag.GetValue())
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := GitModifiersTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_dd_git_modifiers", valueTag.GetKey())
		assert.Equal(t, "a", valueTag.GetValue())
	})

	t.Run("Tag description tests", func(t *testing.T) {
		tag := tags.Tag{}
		defaultDescription := tag.GetDescription()
		cwd, _ := os.Getwd()
		g := TagGroup{}
		g.InitTagGroup(cwd, nil, nil)
		for _, tag := range g.GetTags() {
			assert.NotEqual(t, defaultDescription, tag.GetDescription())
			assert.NotEqual(t, "", tag.GetDescription())
		}
	})

}

func EvaluateTag(t *testing.T, tag tags.ITag, blame gitservice.GitBlame) tags.ITag {
	tag.Init()
	newTag, err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to evaluate tag", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)

	return newTag
}

func EvaluateTagWithPrefix(t *testing.T, tag tags.ITag, blame gitservice.GitBlame, tagPrefix string) tags.ITag {
	tag.Init()
	tag.SetTagPrefix(tagPrefix)
	newTag, err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to evaluate tag with prefix", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)

	return newTag
}
