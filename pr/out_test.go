package pr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/cloudfoundry-community/github-pr-instances-resource/models/fakes"
	"github.com/cloudfoundry-community/github-pr-instances-resource/pr"
	"github.com/cloudfoundry-community/github-pr-instances-resource/test_helpers"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPut(t *testing.T) {

	tests := []struct {
		description string
		source      pr.Source
		version     pr.Version
		parameters  pr.PutParameters
		pullRequest *models.PullRequest
	}{
		{
			description: "put with no parameters does nothing",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters:  pr.PutParameters{},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can set status on a commit",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status: "success",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom context for the status",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status:  "failure",
				Context: "build",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom base context for the status",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status:      "failure",
				BaseContext: "concourse-ci-custom",
				Context:     "build",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom target url for the status",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status:    "failure",
				TargetURL: "https://targeturl.com/concourse",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can provide a custom description for the status",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status:      "failure",
				Description: "Concourse CI build",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can comment on the pull request",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				Number: 1,
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Comment: "comment",
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can delete previous comments made on the pull request",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				Number: 1,
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				DeletePreviousComments: true,
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, []string{}, false, githubv4.PullRequestStateOpen),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.GetPullRequestReturns(tc.pullRequest, nil)

			git := new(fakes.FakeGit)
			git.RevParseReturns("sha", nil)

			dir := test_helpers.CreateTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run get so we have version and metadata for the put request
			// (This is tested in in_test.go)
			getInput := pr.GetRequest{Source: tc.source, Version: tc.version, Params: pr.GetParameters{}}
			_, err := pr.Get(getInput, github, git, dir)
			require.NoError(t, err)

			putInput := pr.PutRequest{Source: tc.source, Params: tc.parameters}
			output, err := pr.Put(putInput, github, dir)

			// Validate output
			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)
			}

			// Validate method calls put on Github.
			if tc.parameters.Status != "" {
				if assert.Equal(t, 1, github.UpdateCommitStatusCallCount()) {
					commit, baseContext, context, status, targetURL, description := github.UpdateCommitStatusArgsForCall(0)
					assert.Equal(t, tc.version.Ref, commit)
					assert.Equal(t, tc.parameters.BaseContext, baseContext)
					assert.Equal(t, tc.parameters.Context, context)
					assert.Equal(t, tc.parameters.TargetURL, targetURL)
					assert.Equal(t, tc.parameters.Description, description)
					assert.Equal(t, tc.parameters.Status, status)
				}
			}

			if tc.parameters.Comment != "" {
				if assert.Equal(t, 1, github.PostCommentCallCount()) {
					pr, comment := github.PostCommentArgsForCall(0)
					assert.Equal(t, tc.pullRequest.Number, pr)
					assert.Equal(t, tc.parameters.Comment, comment)
				}
			}

			if tc.parameters.DeletePreviousComments {
				if assert.Equal(t, 1, github.DeletePreviousCommentsCallCount()) {
					pr := github.DeletePreviousCommentsArgsForCall(0)
					assert.Equal(t, tc.pullRequest.Number, pr)
				}
			}
		})
	}
}

func TestVariableSubstitution(t *testing.T) {

	var (
		variableName  = "BUILD_JOB_NAME"
		variableValue = "my-job"
		variableURL   = "https://concourse-ci.org/"
	)

	tests := []struct {
		description       string
		source            pr.Source
		version           pr.Version
		parameters        pr.PutParameters
		expectedComment   string
		expectedTargetURL string
		pullRequest       *models.PullRequest
	}{

		{
			description: "we can substitute environment variables for Comment",
			source: pr.Source{
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Comment: fmt.Sprintf("$%s", variableName),
			},
			expectedComment: variableValue,
			pullRequest:     test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we can substitute environment variables for TargetURL",
			source: pr.Source{
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Status:    "failure",
				TargetURL: fmt.Sprintf("%s$%s", variableURL, variableName),
			},
			expectedTargetURL: fmt.Sprintf("%s%s", variableURL, variableValue),
			pullRequest:       test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},

		{
			description: "we do not substitute variables other then concourse build metadata",
			source: pr.Source{
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
			},
			version: pr.Version{
				Ref: "commit1",
			},
			parameters: pr.PutParameters{
				Comment: "$THIS_IS_NOT_SUBSTITUTED",
			},
			expectedComment: "$THIS_IS_NOT_SUBSTITUTED",
			pullRequest:     test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.GetPullRequestReturns(tc.pullRequest, nil)

			git := new(fakes.FakeGit)
			git.RevParseReturns("sha", nil)

			dir := test_helpers.CreateTestDirectory(t)
			defer os.RemoveAll(dir)

			// Run get so we have version and metadata for the put request
			getInput := pr.GetRequest{Source: tc.source, Version: tc.version, Params: pr.GetParameters{}}
			_, err := pr.Get(getInput, github, git, dir)
			require.NoError(t, err)

			oldValue := os.Getenv(variableName)
			defer os.Setenv(variableName, oldValue)

			os.Setenv(variableName, variableValue)

			putInput := pr.PutRequest{Source: tc.source, Params: tc.parameters}
			_, err = pr.Put(putInput, github, dir)

			if tc.parameters.TargetURL != "" {
				if assert.Equal(t, 1, github.UpdateCommitStatusCallCount()) {
					_, _, _, _, targetURL, _ := github.UpdateCommitStatusArgsForCall(0)
					assert.Equal(t, tc.expectedTargetURL, targetURL)
				}
			}

			if tc.parameters.Comment != "" {
				if assert.Equal(t, 1, github.PostCommentCallCount()) {
					_, comment := github.PostCommentArgsForCall(0)
					assert.Equal(t, tc.expectedComment, comment)
				}
			}

		})
	}
}
