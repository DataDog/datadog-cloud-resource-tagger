package code2cloud

import (
	"testing"

	"github.com/Datadog/cloud-resource-tagger/src/common/structure"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
	tfStructure "github.com/Datadog/cloud-resource-tagger/src/terraform/structure"
	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	t.Run("ResourceSignatureTagCreation", func(t *testing.T) {
		p := &tfStructure.TerraformParser{}
		p.Init("../../../../src/terraform/structure/", nil)
		defer p.Close()
		filePath := "../../../../src/terraform/structure/main_tagged.tf"

		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}
		for _, block := range parsedBlocks {
			tag := ResourceSignatureTag{}
			valueTag := EvaluateTag(t, &tag, block)
			assert.Equal(t, "dd_git_resource_signature", valueTag.GetKey())
			assert.Equal(t, "module.complete_sg", valueTag.GetValue())
		}

	})
}

func EvaluateTag(t *testing.T, tag tags.ITag, block structure.IBlock) tags.ITag {
	tag.Init()
	newTag, err := tag.CalculateValue(block)
	if err != nil {
		assert.Fail(t, "Failed to evaluate tag", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)

	return newTag
}
