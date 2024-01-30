// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gitservice

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	"github.com/Datadog/cloud-resource-tagger/src/common/structure"
)

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
	GitUserEmail  string
}

func GetPreviousBlameResult(gitSvc *GitService, filePath string) (*git.BlameResult, *object.Commit) {
	if gitSvc.repository == nil {
		return nil, nil
	}
	ref, err := gitSvc.repository.Head()
	if err != nil {
		return nil, nil
	}
	commit, err := gitSvc.repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil
	}
	parentIter := commit.Parents()
	previousCommit, err := parentIter.Next()
	if err != nil {
		return nil, nil
	}

	var previousBlameResult *git.BlameResult
	result, ok := gitSvc.PreviousBlameByFile.Load(filePath)
	if !ok {
		return nil, nil
	}

	previousBlameResult = result.(*git.BlameResult)
	return previousBlameResult, previousCommit
}

func NewGitBlame(relativeFilePath string, filePath string, lines structure.Lines, blameResult *git.BlameResult, gitSvc *GitService) *GitBlame {
	gitBlame := GitBlame{GitOrg: gitSvc.organization, GitRepository: gitSvc.repoName, BlamesByLine: map[int]*git.Line{}, FilePath: relativeFilePath, GitUserEmail: gitSvc.currentUserEmail}
	startLine := lines.Start - 1 // the lines in blameResult.Lines start from zero while the lines range start from 1
	endLine := lines.End - 1
	previousBlameResult, previousCommit := GetPreviousBlameResult(gitSvc, filePath)

	for line := startLine; line <= endLine; line++ {
		if line >= len(blameResult.Lines) {
			logger.Warning(fmt.Sprintf("Index out of bound on parsed file %s", relativeFilePath))
			return &gitBlame
		}
		gitBlame.BlamesByLine[line+1] = blameResult.Lines[line]

		// Check if the line has been removed in the current state of the file
		if previousBlameResult != nil && len(previousBlameResult.Lines) > len(blameResult.Lines) {
			if previousBlameResult.Lines[line].Text != blameResult.Lines[line].Text {
				// The line has been removed, so update the git commit id
				gitBlame.BlamesByLine[line+1].Hash = previousCommit.Hash
			}
		}
	}

	return &gitBlame
}

func (g *GitBlame) GetLatestCommit() (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, v := range g.BlamesByLine {
		if v == nil {
			// This line was added/edited but not committed yet, so latest commit is nil
			return nil
		}
		if latestDate.Before(v.Date) &&
			// Commit was not made by CI, i.e. github actions (for now)
			!strings.Contains(v.Author, "[bot]") && !strings.Contains(v.Author, "github-actions") {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}

func (g *GitBlame) GetLatestCommitForLine(line int) (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for k, v := range g.BlamesByLine {
		// not the relevant line
		if k != line {
			continue
		}

		if v == nil {
			// This line was added/edited but not committed yet, so latest commit is nil
			return nil
		}

		if latestDate.Before(v.Date) &&
			// Commit was not made by CI, i.e. github actions (for now)
			!strings.Contains(v.Author, "[bot]") && !strings.Contains(v.Author, "github-actions") {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}

func GetResourceLineDefinition(blamesByLine map[int]*git.Line) int {
	keys := make([]int, 0)
	for k := range blamesByLine {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys[0]
}

func (g *GitBlame) GetFirstCommit() (firstCommit *git.Line) {
	minimum := time.Time{}

	// to get the relevant commit only for the resource definition line
	resourceLineDef := GetResourceLineDefinition(g.BlamesByLine)
	v := g.BlamesByLine[resourceLineDef]

	if v == nil {
		// This line was added/edited but not committed yet, so latest commit is nil
		return nil
	}
	if (v.Date.Before(minimum) || minimum.IsZero()) &&
		// Commit was not made by CI, i.e. github actions (for now)
		!strings.Contains(v.Author, "[bot]") && !strings.Contains(v.Author, "github-actions") {
		minimum = v.Date
		firstCommit = v
	}
	return
}
