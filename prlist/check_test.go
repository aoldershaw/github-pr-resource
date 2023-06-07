package prlist_test

import (
	"testing"
	"time"

	"github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/cloudfoundry-community/github-pr-instances-resource/models/fakes"
	"github.com/cloudfoundry-community/github-pr-instances-resource/prlist"
	"github.com/cloudfoundry-community/github-pr-instances-resource/test_helpers"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

var (
	testPullRequests = []*models.PullRequest{
		test_helpers.CreateTestPR(1, "master", true, false, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(2, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(3, "master", false, false, 0, nil, true, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(4, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(5, "master", false, true, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(6, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(7, "develop", false, false, 0, []string{"enhancement"}, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(8, "master", false, false, 1, []string{"wontfix"}, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(9, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
		test_helpers.CreateTestPR(10, "master", false, false, 0, nil, false, githubv4.PullRequestStateClosed),
		test_helpers.CreateTestPR(11, "master", false, false, 0, nil, false, githubv4.PullRequestStateMerged),
		test_helpers.CreateTestPR(12, "master", false, false, 0, nil, false, githubv4.PullRequestStateOpen),
	}
)

func TestCheck(t *testing.T) {
	tests := []struct {
		description  string
		source       prlist.Source
		version      *prlist.Version
		files        [][]string
		pullRequests []*models.PullRequest
		expected     prlist.CheckResponse
	}{
		{
			description: "check returns the latest version if there is no previous",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version:      nil,
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[1,2,3,4,5,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check returns the previous version when its still latest",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version:      &prlist.Version{PRs: "[2]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests[1:2],
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[2]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check returns all new versions since the last",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
			},
			version:      &prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[1,2,3,4,5,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check will only return versions that match the specified paths",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				Paths: []string{"terraform/*/*.tf", "terraform/*/*/*.tf"},
			},
			version:      &prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},
				{"terraform/modules/ecs/main.tf", "README.md"},
				{"terraform/modules/variables.tf", "travis.yml"},
			},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[2,3]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check will skip versions which only match the ignore paths",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				IgnorePaths: []string{"*.md", "*.yml"},
			},
			version:      &prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			files: [][]string{
				{"README.md", "travis.yml"},                      // Applies to PR 1
				{"terraform/modules/ecs/main.tf", "README.md"},   // Applies to PR 2
				{"terraform/modules/variables.tf", "travis.yml"}, // Applies to PR 3
				// Subsequent calls will be empty
			},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[2,3]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check correctly ignores [skip ci] when specified",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				DisableCISkip: true,
			},
			version:      &prlist.Version{PRs: "[2]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[2]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[1,2,3,4,5,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check correctly ignores drafts when drafts are ignored",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				IgnoreDrafts: true,
			},
			version:      &prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[1,2,4,5,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check does not ignore drafts when drafts are not ignored",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				IgnoreDrafts: false,
			},
			version:      &prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[4]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[1,2,3,4,5,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check correctly ignores cross repo pull requests",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				DisableForks: true,
			},
			version:      &prlist.Version{PRs: "[6]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[6]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[1,2,3,4,6,7,8,9,12]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check supports specifying base branch",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				BaseBranch: "develop",
			},
			version:      nil,
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.NewVersion([]*models.PullRequest{testPullRequests[6]}),
			},
		},

		{
			description: "check correctly ignores PRs with no approved reviews when specified",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				RequiredReviewApprovals: 1,
			},
			version:      &prlist.Version{PRs: "[9]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[9]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[8]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},

		{
			description: "check returns latest version from a PR with at least one of the desired labels on it",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				Labels: []string{"enhancement"},
			},
			version:      nil,
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.NewVersion([]*models.PullRequest{testPullRequests[6]}),
			},
		},

		{
			description: "check returns latest version from a PR with a single state filter",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				States: []githubv4.PullRequestState{githubv4.PullRequestStateClosed},
			},
			version:      nil,
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.NewVersion([]*models.PullRequest{testPullRequests[9]}),
			},
		},

		{
			description: "check filters out versions from a PR which do not match the state filter",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				States: []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
			},
			version:      nil,
			pullRequests: testPullRequests[9:12],
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.NewVersion([]*models.PullRequest{testPullRequests[11]}),
			},
		},

		{
			description: "check returns versions from a PR with multiple state filters",
			source: prlist.Source{
				GithubConfig: models.GithubConfig{
					Repository: "itsdalmo/test-repository",
				},
				CommonConfig: models.CommonConfig{
					AccessToken: "oauthtoken",
				},
				States: []githubv4.PullRequestState{githubv4.PullRequestStateClosed, githubv4.PullRequestStateMerged},
			},
			version:      &prlist.Version{PRs: "[12]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
			pullRequests: testPullRequests,
			files:        [][]string{},
			expected: prlist.CheckResponse{
				prlist.Version{PRs: "[12]", Timestamp: time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")},
				prlist.Version{PRs: "[10,11]", Timestamp: time.Now().Format("2006-01-02 15:04:05")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			github := new(fakes.FakeGithub)
			pullRequests := []*models.PullRequest{}
			filterStates := []githubv4.PullRequestState{githubv4.PullRequestStateOpen}
			if len(tc.source.States) > 0 {
				filterStates = tc.source.States
			}
			for i := range tc.pullRequests {
				for j := range filterStates {
					if filterStates[j] == tc.pullRequests[i].PullRequestObject.State {
						pullRequests = append(pullRequests, tc.pullRequests[i])
						break
					}
				}
			}
			github.ListPullRequestsReturns(pullRequests, nil)

			for i, file := range tc.files {
				github.ListModifiedFilesReturnsOnCall(i, file, nil)
			}

			input := prlist.CheckRequest{Source: tc.source, Version: tc.version}
			output, err := prlist.Check(input, github)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.expected, output)
			}
			assert.Equal(t, 1, github.ListPullRequestsCallCount())
		})
	}
}

func TestContainsSkipCI(t *testing.T) {
	tests := []struct {
		description string
		message     string
		want        bool
	}{
		{
			description: "does not just match any symbol in the regexp",
			message:     "(",
			want:        false,
		},
		{
			description: "does not match when it should not",
			message:     "test",
			want:        false,
		},
		{
			description: "matches [ci skip]",
			message:     "[ci skip]",
			want:        true,
		},
		{
			description: "matches [skip ci]",
			message:     "[skip ci]",
			want:        true,
		},
		{
			description: "matches trailing skip ci",
			message:     "trailing [skip ci]",
			want:        true,
		},
		{
			description: "matches leading skip ci",
			message:     "[skip ci] leading",
			want:        true,
		},
		{
			description: "is case insensitive",
			message:     "case[Skip CI]insensitive",
			want:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got := prlist.ContainsSkipCI(tc.message)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFilterPath(t *testing.T) {
	cases := []struct {
		description string
		pattern     string
		files       []string
		want        []string
	}{
		{
			description: "returns all matching files",
			pattern:     "*.txt",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"file1.txt",
			},
		},
		{
			description: "works with wildcard",
			pattern:     "test/*",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "excludes unmatched files",
			pattern:     "*/*.txt",
			files: []string{
				"test/file1.go",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "handles prefix matches",
			pattern:     "foo/",
			files: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
				"foo",
				"bar",
				"bar/a.txt",
			},
			want: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := prlist.FilterPath(tc.files, tc.pattern)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFilterIgnorePath(t *testing.T) {
	cases := []struct {
		description string
		pattern     string
		files       []string
		want        []string
	}{
		{
			description: "excludes all matching files",
			pattern:     "*.txt",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"test/file2.txt",
			},
		},
		{
			description: "works with wildcard",
			pattern:     "test/*",
			files: []string{
				"file1.txt",
				"test/file2.txt",
			},
			want: []string{
				"file1.txt",
			},
		},
		{
			description: "includes unmatched files",
			pattern:     "*/*.txt",
			files: []string{
				"test/file1.go",
				"test/file2.txt",
			},
			want: []string{
				"test/file1.go",
			},
		},
		{
			description: "handles prefix matches",
			pattern:     "foo/",
			files: []string{
				"foo/a",
				"foo/a.txt",
				"foo/a/b/c/d.txt",
				"foo",
				"bar",
				"bar/a.txt",
			},
			want: []string{
				"foo",
				"bar",
				"bar/a.txt",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := prlist.FilterIgnorePath(tc.files, tc.pattern)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestIsInsidePath(t *testing.T) {
	cases := []struct {
		description string
		parent      string

		expectChildren    []string
		expectNotChildren []string

		want bool
	}{
		{
			description: "basic test",
			parent:      "foo/bar",
			expectChildren: []string{
				"foo/bar",
				"foo/bar/baz",
			},
			expectNotChildren: []string{
				"foo/barbar",
				"foo/baz/bar",
			},
		},
		{
			description: "does not match parent directories against child files",
			parent:      "foo/",
			expectChildren: []string{
				"foo/bar",
			},
			expectNotChildren: []string{
				"foo",
			},
		},
		{
			description: "matches parents without trailing slash",
			parent:      "foo/bar",
			expectChildren: []string{
				"foo/bar",
				"foo/bar/baz",
			},
		},
		{
			description: "handles children that are shorter than the parent",
			parent:      "foo/bar/baz",
			expectNotChildren: []string{
				"foo",
				"foo/bar",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			for _, expectedChild := range tc.expectChildren {
				if !prlist.IsInsidePath(tc.parent, expectedChild) {
					t.Errorf("Expected \"%s\" to be inside \"%s\"", expectedChild, tc.parent)
				}
			}

			for _, expectedNotChild := range tc.expectNotChildren {
				if prlist.IsInsidePath(tc.parent, expectedNotChild) {
					t.Errorf("Expected \"%s\" to not be inside \"%s\"", expectedNotChild, tc.parent)
				}
			}
		})
	}
}
