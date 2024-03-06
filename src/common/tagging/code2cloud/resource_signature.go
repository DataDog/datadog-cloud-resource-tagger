// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package code2cloud

import (
	"fmt"
	"reflect"

	"github.com/Datadog/cloud-resource-tagger/src/common/structure"
	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type ResourceSignatureTag struct {
	tags.Tag
}

func (t *ResourceSignatureTag) Init() {
	t.Key = "dd_git_resource_signature"
}

func (t *ResourceSignatureTag) CalculateValue(data interface{}) (tags.ITag, error) {
	block, ok := data.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to structure.IBlock, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	hclBlock, ok := block.GetRawBlock().(*hclwrite.Block)
	if !ok {
		return nil, fmt.Errorf("failed to convert block to *hclwrite.Block, which is required to calculte tag value. Type of block: %s", reflect.TypeOf(block.GetRawBlock()))
	}

	resourceSignature := hclBlock.Type() + "." + block.GetResourceID()
	fmt.Printf("ResourceSignatureTag: %v\n", resourceSignature)
	return &tags.Tag{Key: t.Key, Value: resourceSignature}, nil
}

func (t *ResourceSignatureTag) GetDescription() string {
	return "Signature of the resource defined in the file"
}
