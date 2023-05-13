package prlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	resource "github.com/aoldershaw/github-pr-resource"
	"github.com/shurcooL/githubv4"
)

// Source represents the configuration for the resource.
type Source struct {
	resource.CommonConfig
	resource.GithubConfig
	Paths                   []string                    `json:"paths"`
	IgnorePaths             []string                    `json:"ignore_paths"`
	DisableCISkip           bool                        `json:"disable_ci_skip"`
	DisableForks            bool                        `json:"disable_forks"`
	IgnoreDrafts            bool                        `json:"ignore_drafts"`
	BaseBranch              string                      `json:"base_branch"`
	RequiredReviewApprovals int                         `json:"required_review_approvals"`
	PRNumbersHigherThan     int                         `json:"pr_numbers_higher_than"`
	Labels                  []string                    `json:"labels"`
	States                  []githubv4.PullRequestState `json:"states"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AccessToken == "" {
		return errors.New("access_token must be set")
	}
	if s.Repository == "" {
		return errors.New("repository must be set")
	}
	if s.V3Endpoint != "" && s.V4Endpoint == "" {
		return errors.New("v4_endpoint must be set together with v3_endpoint")
	}
	if s.V4Endpoint != "" && s.V3Endpoint == "" {
		return errors.New("v3_endpoint must be set together with v4_endpoint")
	}
	for _, state := range s.States {
		switch state {
		case githubv4.PullRequestStateOpen:
		case githubv4.PullRequestStateClosed:
		case githubv4.PullRequestStateMerged:
		default:
			return errors.New(fmt.Sprintf("states value \"%s\" must be one of: OPEN, MERGED, CLOSED", state))
		}
	}
	return nil
}

type Version struct {
	// JSON encoded list of PR numbers.
	PRs string `json:"prs"`
	// Time when the version was initially generated.
	Timestamp string `json:"timestamp"`
}

// NewVersion constructs a new Version.
func NewVersion(prs []*resource.PullRequest) Version {
	numbers := make([]int, len(prs))
	for i, pr := range prs {
		numbers[i] = pr.Number
	}
	data, err := json.Marshal(numbers)
	if err != nil {
		panic(err)
	}
	return Version{
		PRs:       string(data),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// PRData represents a single PR in the get response file.
type PRData struct {
	Number int `json:"number"`
}
