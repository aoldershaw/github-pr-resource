package pr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	resource "github.com/aoldershaw/github-pr-resource"
)

func Get(request GetRequest, github resource.Github, git resource.Git, outputDir string) (*GetResponse, error) {
	if request.Params.SkipDownload {
		return &GetResponse{Version: request.Version}, nil
	}

	pull, err := github.GetPullRequest(request.Source.Number, request.Version.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pull request: %s", err)
	}

	// Initialize and pull the base for the PR
	if err := git.Init(&pull.BaseRefName); err != nil {
		return nil, err
	}
	if err := git.Pull(pull.Repository.URL, pull.BaseRefName, request.Params.GitDepth, request.Params.Submodules, request.Params.FetchTags); err != nil {
		return nil, err
	}

	// Get the last commit SHA in base for the metadata
	baseSHA, err := git.RevParse(pull.BaseRefName)
	if err != nil {
		return nil, err
	}

	// Fetch the PR and merge the specified commit into the base
	if err := git.Fetch(pull.Repository.URL, pull.Number, request.Params.GitDepth, request.Params.Submodules, false); err != nil {
		return nil, err
	}

	// Create the metadata
	var metadata resource.Metadata
	metadata.Add("pr", strconv.Itoa(pull.Number))
	metadata.Add("title", pull.Title)
	metadata.Add("url", pull.URL)
	metadata.Add("head_name", pull.HeadRefName)
	metadata.Add("head_sha", pull.Tip.OID)
	metadata.Add("base_name", pull.BaseRefName)
	metadata.Add("base_sha", baseSHA)
	metadata.Add("message", pull.Tip.Message)
	metadata.Add("author", pull.Tip.Author.User.Login)
	metadata.Add("author_email", pull.Tip.Author.Email)
	metadata.Add("state", string(pull.State))

	// Write version and metadata for reuse in PUT
	path := filepath.Join(outputDir, ".git", "resource")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %s", err)
	}
	b, err := json.Marshal(request.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(path, "version.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write version: %s", err)
	}
	b, err = json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(path, "metadata.json"), b, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %s", err)
	}

	for _, d := range metadata {
		filename := d.Name
		content := []byte(d.Value)
		if err := ioutil.WriteFile(filepath.Join(path, filename), content, 0644); err != nil {
			return nil, fmt.Errorf("failed to write metadata file %s: %s", filename, err)
		}
	}

	switch tool := request.Params.IntegrationTool; tool {
	case "rebase":
		if err := git.Rebase(pull.BaseRefName, pull.Tip.OID, request.Params.Submodules); err != nil {
			return nil, err
		}
	case "merge", "":
		if err := git.Merge(pull.Tip.OID, request.Params.Submodules); err != nil {
			return nil, err
		}
	case "checkout":
		if err := git.Checkout(pull.HeadRefName, pull.Tip.OID, request.Params.Submodules); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid integration tool specified: %s", tool)
	}

	if request.Source.GitCryptKey != "" {
		if err := git.GitCryptUnlock(request.Source.GitCryptKey); err != nil {
			return nil, err
		}
	}

	if request.Params.ListChangedFiles {
		changedFiles, err := github.ListModifiedFiles(request.Source.Number)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch list of changed files: %s", err)
		}

		changedFilesBytes := []byte(strings.Join(changedFiles, "\n") + "\n")
		if err := ioutil.WriteFile(filepath.Join(path, "changed_files"), changedFilesBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to write file list: %s", err)
		}
	}

	return &GetResponse{
		Version:  request.Version,
		Metadata: metadata,
	}, nil
}

type GetParameters struct {
	SkipDownload     bool   `json:"skip_download"`
	IntegrationTool  string `json:"integration_tool"`
	GitDepth         int    `json:"git_depth"`
	Submodules       bool   `json:"submodules"`
	ListChangedFiles bool   `json:"list_changed_files"`
	FetchTags        bool   `json:"fetch_tags"`
}

type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

type GetResponse struct {
	Version  Version           `json:"version"`
	Metadata resource.Metadata `json:"metadata,omitempty"`
}
