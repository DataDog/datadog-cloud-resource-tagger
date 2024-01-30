// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package code2cloud

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Datadog/cloud-resource-tagger/src/common/tagging/tags"
)

type TraceTag struct {
	tags.Tag
}

func (t *TraceTag) Init() {
	t.Key = tags.TraceTagKey
}

func (t *TraceTag) CalculateValue(_ interface{}) (tags.ITag, error) {
	uuidv4, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create a new uuidv4")
	}
	return &tags.Tag{Key: t.Key, Value: uuidv4.String()}, nil
}

func (t *TraceTag) GetDescription() string {
	return "A UUID tag that allows easily finding the root IaC config of the resource"
}
