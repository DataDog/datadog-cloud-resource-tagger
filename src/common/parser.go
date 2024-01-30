// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package common

import "github.com/Datadog/cloud-resource-tagger/src/common/structure"

type IParser interface {
	Init(rootDir string, args map[string]string)
	Name() string
	ValidFile(filePath string) bool
	ParseFile(filePath string) ([]structure.IBlock, error)
	WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetSupportedFileExtensions() []string
	Close()
}
