package utils

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff-index", "--name-only", "HEAD")
	fileCommandOutput, err := cmd.Output()
	changedFiles := ParseChangedFileList(string(fileCommandOutput))
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	return changedFiles, nil
}

func ParseChangedFileList(files string) []string {
	result := []string{}
	scanner := bufio.NewScanner(strings.NewReader(files))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	return result
}
