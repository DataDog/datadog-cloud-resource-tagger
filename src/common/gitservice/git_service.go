// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.
package gitservice

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/Datadog/cloud-resource-tagger/src/common/logger"
	"github.com/Datadog/cloud-resource-tagger/src/common/structure"
)

type GitService struct {
	gitRootDir          string
	scanPathFromRoot    string
	repository          *git.Repository
	remoteURL           string
	organization        string
	repoName            string
	BlameByFile         *sync.Map
	PreviousBlameByFile *sync.Map
	currentUserEmail    string
	repositoryUrl       string
}

var gitGraphLock sync.Mutex

func NewGitService(rootDir string) (*GitService, error) {
	var repository *git.Repository
	var err error
	rootDirIter, _ := filepath.Abs(rootDir)
	for {
		repository, err = git.PlainOpen(".")
		if err == nil {
			break
		}
		newRootDir := filepath.Dir(rootDirIter)
		if rootDirIter == newRootDir {
			break
		}
		rootDirIter = newRootDir
	}
	if err != nil {
		return nil, err
	}

	scanAbsDir, _ := filepath.Abs(rootDir)
	scanPathFromRoot, _ := filepath.Rel(rootDirIter, scanAbsDir)

	gitService := GitService{
		gitRootDir:          rootDir,
		scanPathFromRoot:    scanPathFromRoot,
		repository:          repository,
		BlameByFile:         &sync.Map{},
		PreviousBlameByFile: &sync.Map{},
	}
	err = gitService.setOrgAndName()
	gitService.currentUserEmail = GetGitUserEmail()

	return &gitService, err
}

func (g *GitService) setOrgAndName() error {
	// get remotes to find the repository's url
	remotes, err := g.repository.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes, err: %s", err)
	}

	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			g.remoteURL = remote.Config().URLs[0]

			g.repositoryUrl = strings.TrimSuffix(g.remoteURL, ".git")

			// get endpoint structured like '/github.com/Datadog/cloud-resource-tagger.git
			endpoint, err := transport.NewEndpoint(g.remoteURL)
			if err != nil {
				return err
			}
			// remove leading '/' from path and trailing '.git. suffix, then split by '/'
			endpointPathParts := strings.Split(strings.TrimSuffix(strings.TrimLeft(endpoint.Path, "/"), ".git"), "/")
			if len(endpointPathParts) < 2 {
				return fmt.Errorf("invalid format of endpoint path: %s", endpoint.Path)
			}
			g.organization = endpointPathParts[0]
			g.repoName = strings.Join(endpointPathParts[1:], "/")
			break
		}
	}

	return nil
}

func (g *GitService) ComputeRelativeFilePath(fp string) string {
	if strings.HasPrefix(fp, g.gitRootDir) {
		res, _ := filepath.Rel(g.gitRootDir, fp)
		return filepath.Join(g.scanPathFromRoot, res)
	}
	scanPathIter := g.scanPathFromRoot
	parent := filepath.Dir(fp)
	for {
		_, child := filepath.Split(scanPathIter)
		if parent != child {
			break
		}
		scanPathIter, _ = filepath.Split(scanPathIter)
	}
	return filepath.Join(scanPathIter, fp)
}

func (g *GitService) GetBlameForFileLines(filePath string, lines structure.Lines) (*GitBlame, error) {
	logger.Info(fmt.Sprintf("Getting git blame for %v (%v:%v)", filePath, lines.Start, lines.End))
	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	blame, ok := g.BlameByFile.Load(filePath)
	if ok {
		return NewGitBlame(relativeFilePath, filePath, lines, blame.(*git.BlameResult), g), nil
	}

	var err error
	blame, err = g.GetFileBlame(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}

	g.BlameByFile.Store(filePath, blame)

	return NewGitBlame(relativeFilePath, filePath, lines, blame.(*git.BlameResult), g), nil
}

func (g *GitService) GetOrganization() string {
	return g.organization
}

func (g *GitService) GetRepoName() string {
	return g.repoName
}

func (g *GitService) CommitChanges(filePath string) error {
	gitGraphLock.Lock() // Git is a graph, different files can lead to graph scans interfering with each other
	defer gitGraphLock.Unlock()
	w, err := g.repository.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get repository worktree for file %s because of error %s", filePath, err)
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		return nil
	}

	for path, _ := range status {
		_, err := w.Add(path)
		if err != nil {
			return fmt.Errorf("failed to add file %s to git index because of error %s", path, err)
		}
	}

	_, err = w.Commit("Adding tags from datadog-cloud-resource-tagger", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Datadog Cloud Resource Tagger",
			Email: g.currentUserEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes for file %s because of error %s", filePath, err)
	}
	return nil
}

func (g *GitService) GetFileBlame(filePath string) (*git.BlameResult, error) {
	blame, ok := g.BlameByFile.Load(filePath)
	if ok {
		return blame.(*git.BlameResult), nil
	}

	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	var selectedCommit *object.Commit

	gitGraphLock.Lock() // Git is a graph, different files can lead to graph scans interfering with each other
	defer gitGraphLock.Unlock()
	head, err := g.repository.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository HEAD for file %s because of error %s", filePath, err)
	}
	selectedCommit, err = g.repository.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to find commit %s ", head.Hash().String())
	}

	parentIter := selectedCommit.Parents()
	previousCommit, err := parentIter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get previous commit: %s", err)
	}
	blame, err = git.Blame(selectedCommit, relativeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}
	previousBlame, err := git.Blame(previousCommit, relativeFilePath)
	if err != nil {
		fmt.Printf("failed to get previous blame for previous commit of file %s because of error %s", filePath, err)
	}
	if previousBlame != nil {
		g.PreviousBlameByFile.Store(filePath, previousBlame)
	}
	g.BlameByFile.Store(filePath, blame)

	return blame.(*git.BlameResult), nil
}

func GetGitUserEmail() string {
	log.SetOutput(io.Discard)
	cmd := exec.Command("git", "config", "user.email")
	email, err := cmd.Output()
	stdout := os.Stdout
	log.SetOutput(stdout)
	if err != nil {
		logger.Debug(fmt.Sprintf("unable to get current git user email: %s", err))
		return ""
	}
	return strings.ReplaceAll(string(email), "\n", "")
}
