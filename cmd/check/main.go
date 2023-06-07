package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	models "github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/cloudfoundry-community/github-pr-instances-resource/pr"
	"github.com/cloudfoundry-community/github-pr-instances-resource/prlist"
)

type Request struct {
	Source struct {
		Number int `json:"number"`
	} `json:"source"`
}

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read stdin: %v", err)
	}

	var request Request
	if err := json.Unmarshal(stdin, &request); err != nil {
		log.Fatalf("failed to unmarshal request: %v", err)
	}

	if request.Source.Number == 0 {
		checkPRList(stdin)
	} else {
		checkPR(stdin)
	}
}

func checkPRList(stdin []byte) {
	decoder := json.NewDecoder(bytes.NewReader(stdin))
	decoder.DisallowUnknownFields()

	var request prlist.CheckRequest
	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %v", err)
	}

	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %v", err)
	}
	github, err := models.NewGithubClient(request.Source.CommonConfig, request.Source.GithubConfig)
	if err != nil {
		log.Fatalf("failed to create github manager: %v", err)
	}
	response, err := prlist.Check(request, github)
	if err != nil {
		log.Fatalf("check failed: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
}

func checkPR(stdin []byte) {
	decoder := json.NewDecoder(bytes.NewReader(stdin))
	decoder.DisallowUnknownFields()

	var request pr.CheckRequest
	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %v", err)
	}

	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %v", err)
	}
	const repoDir = "/tmp/git-repo"
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}

	// We never need git-lfs when we check for new commits, so always disable it.
	git, err := models.NewGitClient(request.Source.CommonConfig, true, repoDir, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create git manager: %v", err)
	}
	response, err := pr.Check(request, git)
	if err != nil {
		log.Fatalf("check failed: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
}
