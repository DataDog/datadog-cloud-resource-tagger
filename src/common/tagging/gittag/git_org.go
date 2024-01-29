package gittag

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/gitservice"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type GitOrgTag struct {
	tags.Tag
}

func (t *GitOrgTag) Init() {
	t.Key = "dd_git_org"
}

func (t *GitOrgTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &tags.Tag{Key: t.Key, Value: gitBlame.GitOrg}, nil
}

func (t *GitOrgTag) GetDescription() string {
	return "The entity which owns the repository where this resource is provisioned in IaC"
}
