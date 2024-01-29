package gittag

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitCreatedAtTag struct {
	tags.Tag
}

func (t *GitCreatedAtTag) Init() {
	t.Key = tags.GitCreatedAtTagKey
}

func (t *GitCreatedAtTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	firstCommit := gitBlame.GetFirstCommit()

	if firstCommit == nil {
		return nil, fmt.Errorf("first commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: firstCommit.Date.UTC().Format("2006-01-02 15:04:05")}, nil
}

func (t *GitCreatedAtTag) GetDescription() string {
	return "When this resource's configuration was first created"
}
