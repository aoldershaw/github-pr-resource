package pr_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/cloudfoundry-community/github-pr-instances-resource/models/fakes"
	"github.com/cloudfoundry-community/github-pr-instances-resource/pr"
	"github.com/cloudfoundry-community/github-pr-instances-resource/test_helpers"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	tests := []struct {
		description    string
		source         pr.Source
		version        pr.Version
		parameters     pr.GetParameters
		pullRequest    *models.PullRequest
		versionString  string
		metadataString string
		files          []models.ChangedFileObject
		filesString    string
	}{
		{
			description: "get works",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters:     pr.GetParameters{},
			pullRequest:    test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
		},
		{
			description: "get supports unlocking with git crypt",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				GitCryptKey: "gitcryptkey",
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters:     pr.GetParameters{},
			pullRequest:    test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
		},
		{
			description: "get supports rebasing",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters: pr.GetParameters{
				IntegrationTool: "rebase",
			},
			pullRequest:    test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
		},
		{
			description: "get supports checkout",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters: pr.GetParameters{
				IntegrationTool: "checkout",
			},
			pullRequest:    test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
		},
		{
			description: "get supports git_depth",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters: pr.GetParameters{
				GitDepth: 2,
			},
			pullRequest:    test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
		},
		{
			description: "get supports list_changed_files",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			parameters: pr.GetParameters{
				ListChangedFiles: true,
			},
			pullRequest: test_helpers.CreateTestPR(1, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
			// files: []models.ChangedFileObject{
			// 	{
			// 		Path: "README.md",
			// 	},
			// 	{
			// 		Path: "Other.md",
			// 	},
			// },
			versionString:  `{"pr":"pr1","commit":"commit1","committed":"0001-01-01T00:00:00Z","approved_review_count":"0","state":"OPEN"}`,
			metadataString: `[{"name":"pr","value":"1"},{"name":"title","value":"pr1 title"},{"name":"url","value":"pr1 url"},{"name":"head_name","value":"pr1"},{"name":"head_sha","value":"oid1"},{"name":"base_name","value":"master"},{"name":"base_sha","value":"sha"},{"name":"message","value":"commit message1"},{"name":"author","value":"login1"},{"name":"author_email","value":"user@example.com"},{"name":"state","value":"OPEN"}]`,
			filesString:    "README.md\nOther.md\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			github.GetPullRequestReturns(tc.pullRequest, nil)

			// if tc.files != nil {
			// 	github.GetChangedFilesReturns(tc.files, nil)
			// }

			git := new(fakes.FakeGit)
			git.RevParseReturns("sha", nil)

			dir := test_helpers.CreateTestDirectory(t)
			defer os.RemoveAll(dir)

			input := pr.GetRequest{Source: tc.source, Version: tc.version, Params: tc.parameters}
			output, err := pr.Get(input, github, git, dir)

			// Validate output
			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)

				// Verify written files
				version := test_helpers.ReadTestFile(t, filepath.Join(dir, ".git", "resource", "version.json"))
				assert.Equal(t, fmt.Sprintf("{\"ref\":\"%s\"}", tc.version.Ref), version)

				metadata := test_helpers.ReadTestFile(t, filepath.Join(dir, ".git", "resource", "metadata.json"))
				assert.Equal(t, tc.metadataString, metadata)

				// Verify individual files
				files := map[string]string{
					"pr":           "1",
					"url":          "pr1 url",
					"head_name":    "pr1",
					"head_sha":     "oid1",
					"base_name":    "master",
					"base_sha":     "sha",
					"message":      "commit message1",
					"author":       "login1",
					"author_email": "user@example.com",
					"title":        "pr1 title",
				}

				for filename, expected := range files {
					actual := test_helpers.ReadTestFile(t, filepath.Join(dir, ".git", "resource", filename))
					assert.Equal(t, expected, actual)
				}

				if tc.files != nil {
					changedFiles := test_helpers.ReadTestFile(t, filepath.Join(dir, ".git", "resource", "changed_files"))
					assert.Equal(t, tc.filesString, changedFiles)
				}
			}

			// Validate Github calls
			if assert.Equal(t, 1, github.GetPullRequestCallCount()) {
				pr, _ := github.GetPullRequestArgsForCall(0)
				assert.Equal(t, tc.source.Number, pr)
			}

			// Validate Git calls
			if assert.Equal(t, 1, git.InitCallCount()) {
				base := git.InitArgsForCall(0)
				assert.Equal(t, tc.pullRequest.BaseRefName, *base)
			}

			if assert.Equal(t, 1, git.PullCallCount()) {
				url, base, depth, submodules, fetchTags := git.PullArgsForCall(0)
				assert.Equal(t, tc.pullRequest.Repository.URL, url)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
				assert.Equal(t, tc.parameters.GitDepth, depth)
				assert.Equal(t, tc.parameters.Submodules, submodules)
				assert.Equal(t, tc.parameters.FetchTags, fetchTags)
			}

			if assert.Equal(t, 1, git.RevParseCallCount()) {
				base := git.RevParseArgsForCall(0)
				assert.Equal(t, tc.pullRequest.BaseRefName, base)
			}

			if assert.Equal(t, 1, git.FetchCallCount()) {
				url, pr, depth, submodules, checkoutBool := git.FetchArgsForCall(0)
				assert.Equal(t, tc.pullRequest.Repository.URL, url)
				assert.Equal(t, tc.pullRequest.Number, pr)
				assert.Equal(t, tc.parameters.GitDepth, depth)
				assert.Equal(t, tc.parameters.Submodules, submodules)
				assert.Equal(t, false, checkoutBool)
			}

			switch tc.parameters.IntegrationTool {
			case "rebase":
				if assert.Equal(t, 1, git.RebaseCallCount()) {
					branch, tip, submodules := git.RebaseArgsForCall(0)
					assert.Equal(t, tc.pullRequest.BaseRefName, branch)
					assert.Equal(t, tc.pullRequest.Tip.OID, tip)
					assert.Equal(t, tc.parameters.Submodules, submodules)
				}
			case "checkout":
				if assert.Equal(t, 1, git.CheckoutCallCount()) {
					branch, sha, submodules := git.CheckoutArgsForCall(0)
					assert.Equal(t, tc.pullRequest.HeadRefName, branch)
					assert.Equal(t, tc.pullRequest.Tip.OID, sha)
					assert.Equal(t, tc.parameters.Submodules, submodules)
				}
			default:
				if assert.Equal(t, 1, git.MergeCallCount()) {
					tip, submodules := git.MergeArgsForCall(0)
					assert.Equal(t, tc.pullRequest.Tip.OID, tip)
					assert.Equal(t, tc.parameters.Submodules, submodules)
				}
			}
			//FIXME
			if tc.source.GitCryptKey != "" {
				if assert.Equal(t, 1, git.GitCryptUnlockCallCount()) {
					key := git.GitCryptUnlockArgsForCall(0)
					assert.Equal(t, tc.source.GitCryptKey, key)
				}
			}
		})
	}
}

func TestGetSkipDownload(t *testing.T) {
	tests := []struct {
		description string
		source      pr.Source
		version     pr.Version
		params      pr.GetParameters
		pullRequest *models.PullRequest
	}{
		{
			description: "skip download works",
			source: pr.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version: pr.Version{
				Ref: "some-ref",
			},
			params: pr.GetParameters{
				SkipDownload: true,
			},
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

			// Run the get and check output
			input := pr.GetRequest{Source: tc.source, Version: tc.version, Params: tc.params}
			output, err := pr.Get(input, github, git, dir)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.version, output.Version)
			}
		})
	}
}
