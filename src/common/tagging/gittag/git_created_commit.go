package gittag

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitCreateCommitTag struct {
	tags.Tag
}

func (t *GitCreateCommitTag) Init() {
	t.Key = "dd_git_create_commit"
}

func (t *GitCreateCommitTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	firstCommit := gitBlame.GetFirstCommit()
	if firstCommit == nil || firstCommit.Hash.IsZero() {
		return &tags.Tag{Key: t.Key, Value: CommitUnavailable}, nil
	}
	return &tags.Tag{Key: t.Key, Value: firstCommit.Hash.String()}, nil
}

func (t *GitCreateCommitTag) GetDescription() string {
	return "The hash of the latest commit which edited this resource"
}
