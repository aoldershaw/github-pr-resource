package pr

import "github.com/cloudfoundry-community/github-pr-instances-resource/models"

func Check(request CheckRequest, git models.Git) (CheckResponse, error) {
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

	return response, nil
}

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type CheckResponse []Version
