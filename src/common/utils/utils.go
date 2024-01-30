// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package utils

import (
	"regexp"
	"strings"
)

// RemoveGcpInvalidChars Source of regex: https://cloud.google.com/compute/docs/labeling-resources
var RemoveGcpInvalidChars = regexp.MustCompile(`[^\p{Ll}\p{Lo}\p{N}_-]`)

func GetLinesFromBytes(bytes []byte) []string {
	return strings.Split(string(bytes), "\n")
}

func SplitStringByComma(input []string) []string {
	var ans []string
	for _, i := range input {
		if strings.Contains(i, ",") {
			ans = append(ans, strings.Split(i, ",")...)
		} else {
			ans = append(ans, i)
		}
	}
	return ans
}

func InSlice[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func SliceInSlices[T comparable](elems [][]T, vSlice []T) bool {
	for _, elemSlice := range elems {
		curSize := len(elemSlice)
		if curSize != len(vSlice) {
			continue
		}
		equalNum := 0
		for i, elem := range elemSlice {
			if elem == vSlice[i] {
				equalNum++
			}
		}
		if equalNum == curSize {
			return true
		}
	}
	return false
}

func FindSubMatchByGroup(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
