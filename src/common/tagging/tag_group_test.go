package tagging

import (
	"testing"

	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

func TestTagGroup(t *testing.T) {
	t.Run("Test tagGroup skip single tag", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"dd_correlation_id"}}
		tagGroup.SetTags([]tags.ITag{&tags.Tag{Key: "dd_correlation_id"}, &tags.Tag{Key: "git_modifiers"}})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.NotEqual(t, "dd_correlation_id", tgs[0].GetKey())
	})

	t.Run("Test tagGroup skip regex tags", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*"}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "dd_correlation_id"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.Equal(t, "dd_correlation_id", tgs[0].GetKey())
	})

	t.Run("Test tagGroup skip multi", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*", "dd_correlation_id"}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "dd_correlation_id"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 0, len(tgs))
	})

	t.Run("Test tag prefix not broke tagGroup skip multi", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*", "dd_correlation_id"}, Options: InitTagGroupOptions{
			TagPrefix: "prefix_",
		}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "dd_correlation_id"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 0, len(tgs))
	})
}
