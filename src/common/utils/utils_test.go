package utils

import (
	"testing"

	"github.com/likexian/gokit/assert"
)

func TestDetermineTopLevelDirectory(t *testing.T) {
	t.Run("TopLevelPath", func(t *testing.T) {
		files := []string{"main.tf", "../a/b/main.tf", "a/b/c/main.tf"}
		topLevelPath := DetermineTopLevelDirectory(files)
		assert.Equal(t, topLevelPath, "main.tf")
	})

	t.Run("NestedPaths", func(t *testing.T) {
		files := []string{"a/b/main.tf", "a/b/c/main.tf"}
		topLevelPath := DetermineTopLevelDirectory(files)
		assert.Equal(t, topLevelPath, "a/b/main.tf")
	})

}
