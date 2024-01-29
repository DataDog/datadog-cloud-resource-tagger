package gittag

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

const CommitUnavailable = "N/A"

type GitModifiedCommitTag struct {
	tags.Tag
}

func (t *GitModifiedCommitTag) Init() {
	t.Key = "dd_git_modified_commit"
}

func (t *GitModifiedCommitTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	latestCommit := gitBlame.GetLatestCommit()
	if latestCommit == nil || latestCommit.Hash.IsZero() {
		return &tags.Tag{Key: t.Key, Value: CommitUnavailable}, nil
	}
	return &tags.Tag{Key: t.Key, Value: latestCommit.Hash.String()}, nil
}

func (t *GitModifiedCommitTag) GetDescription() string {
	return "The hash of the latest commit which edited this resource"
}
