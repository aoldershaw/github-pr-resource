package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	models "github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/cloudfoundry-community/github-pr-instances-resource/pr"
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
	sourceDir := os.Args[1]

	if request.Source.Number == 0 {
		log.Fatalf("can only put when source.number is specified")
	} else {
		putPR(stdin, sourceDir)
	}
}

func putPR(stdin []byte, sourceDir string) {
	decoder := json.NewDecoder(bytes.NewReader(stdin))
	decoder.DisallowUnknownFields()

	var request pr.PutRequest
	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %s", err)
	}
	github, err := models.NewGithubClient(request.Source.CommonConfig, request.Source.GithubConfig)
	if err != nil {
		log.Fatalf("failed to create github manager: %s", err)
	}
	response, err := pr.Put(request, github, sourceDir)
	if err != nil {
		log.Fatalf("put failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}
