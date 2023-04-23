package pr

import (
	"fmt"

	resource "github.com/aoldershaw/github-pr-resource"
)

func Check(request CheckRequest, git resource.Git, manager resource.Github) (CheckResponse, error) {
	// First check commits
	if err := git.Init(nil); err != nil {
		return CheckResponse{}, err
	}

	url := request.Source.RepositoryURL()
	if err := git.Fetch(url, request.Source.Number, 0, false, true); err != nil {
		return nil, err
	}

	var fromCommit *string
	if request.Version != nil {
		fromCommit = &request.Version.Ref
	}
	commits, err := git.RevList(fromCommit, request.Source.Paths, request.Source.IgnorePaths, request.Source.DisableCISkip)
	if err != nil {
		return nil, err
	}

	response := make(CheckResponse, len(commits))
	for i, commit := range commits {
		response[i] = Version{Ref: commit}
	}

	if len(commits) != 0 && len(request.Source.Labels) != 0 {
		// Filter by PR labels
		pull, err := manager.GetPullRequest(request.Source.Number, commits[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get pull request: %s", err)
		}
		labelFound := false
	LabelLoop:
		for _, wantedLabel := range request.Source.Labels {
			for _, targetLabel := range pull.Labels {
				if targetLabel.Name == wantedLabel {
					labelFound = true
					break LabelLoop
				}
			}
		}
		if !labelFound {
			return make(CheckResponse, 0), nil
		}
	}

	return response, nil
}

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type CheckResponse []Version
