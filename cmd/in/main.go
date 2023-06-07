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

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	outputDir := os.Args[1]

	if request.Source.Number == 0 {
		getPRList(stdin, outputDir)
	} else {
		getPR(stdin, outputDir)
	}
}

func getPRList(stdin []byte, outputDir string) {
	decoder := json.NewDecoder(bytes.NewReader(stdin))
	decoder.DisallowUnknownFields()

	var request prlist.GetRequest
	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %v", err)
	}

	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %v", err)
	}
	response, err := prlist.Get(request, outputDir)
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
}

func getPR(stdin []byte, outputDir string) {
	decoder := json.NewDecoder(bytes.NewReader(stdin))
	decoder.DisallowUnknownFields()

	var request pr.GetRequest
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

	git, err := models.NewGitClient(request.Source.CommonConfig, request.Source.DisableGitLFS, outputDir, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create git manager: %v", err)
	}

	response, err := pr.Get(request, github, git, outputDir)
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
}
