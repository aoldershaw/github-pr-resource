package resource

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type CommonConfig struct {
	AccessToken         string `json:"access_token"`
	SkipSSLVerification bool   `json:"skip_ssl_verification"`
}

type GithubConfig struct {
	Repository      string `json:"repository"`
	HostingEndpoint string `json:"hosting_endpoint"`
	V3Endpoint      string `json:"v3_endpoint"`
	V4Endpoint      string `json:"v4_endpoint"`
}

func (config GithubConfig) RepositoryURL() string {
	hostingEndpoint := config.HostingEndpoint
	if hostingEndpoint == "" {
		hostingEndpoint = "https://github.com"
	}
	return strings.TrimRight(hostingEndpoint, "/") + "/" + config.Repository
}

// Github for testing purposes.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/fake_github.go . Github
type Github interface {
	ListPullRequests([]githubv4.PullRequestState) ([]*PullRequest, error)
	GetPullRequest(int, string) (*PullRequest, error)
	ListModifiedFiles(int) ([]string, error)
	PostComment(int, string) error
	UpdateCommitStatus(string, string, string, string, string, string) error
	DeletePreviousComments(int) error
}

// GithubClient for handling requests to the Github V3 and V4 APIs.
type GithubClient struct {
	HostingEndpoint string
	V3              *github.Client
	V4              *githubv4.Client
	Repository      string
	Owner           string
}

// NewGithubClient ...
func NewGithubClient(common CommonConfig, config GithubConfig) (*GithubClient, error) {
	owner, repository, err := parseRepository(config.Repository)
	if err != nil {
		return nil, err
	}

	// Skip SSL verification for self-signed certificates
	// source: https://github.com/google/go-github/pull/598#issuecomment-333039238
	var ctx context.Context
	if common.SkipSSLVerification {
		insecureClient := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		}
		ctx = context.WithValue(context.TODO(), oauth2.HTTPClient, insecureClient)
	} else {
		ctx = context.TODO()
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: common.AccessToken},
	))

	var v3 *github.Client
	if config.V3Endpoint != "" {
		endpoint, err := url.Parse(config.V3Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse v3 endpoint: %config", err)
		}
		v3, err = github.NewEnterpriseClient(endpoint.String(), endpoint.String(), client)
		if err != nil {
			return nil, err
		}
	} else {
		v3 = github.NewClient(client)
	}

	var v4 *githubv4.Client
	if config.V4Endpoint != "" {
		endpoint, err := url.Parse(config.V4Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse v4 endpoint: %config", err)
		}
		v4 = githubv4.NewEnterpriseClient(endpoint.String(), client)
		if err != nil {
			return nil, err
		}
	} else {
		v4 = githubv4.NewClient(client)
	}

	return &GithubClient{
		HostingEndpoint: config.HostingEndpoint,
		V3:              v3,
		V4:              v4,
		Owner:           owner,
		Repository:      repository,
	}, nil
}

// ListPullRequests gets the last commit on all pull requests with the matching state.
func (m *GithubClient) ListPullRequests(prStates []githubv4.PullRequestState) ([]*PullRequest, error) {
	var query struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						PullRequestObject
						Reviews struct {
							TotalCount int
						} `graphql:"reviews(states: $prReviewStates)"`
						Commits struct {
							Edges []struct {
								Node struct {
									Commit CommitObject
								}
							}
						} `graphql:"commits(last:$commitsLast)"`
						Labels struct {
							Edges []struct {
								Node struct {
									LabelObject
								}
							}
						} `graphql:"labels(first:$labelsFirst)"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"pullRequests(first:$prFirst,states:$prStates,after:$prCursor)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	vars := map[string]interface{}{
		"repositoryOwner": githubv4.String(m.Owner),
		"repositoryName":  githubv4.String(m.Repository),
		"prFirst":         githubv4.Int(100),
		"prStates":        prStates,
		"prCursor":        (*githubv4.String)(nil),
		"commitsLast":     githubv4.Int(1),
		"prReviewStates":  []githubv4.PullRequestReviewState{githubv4.PullRequestReviewStateApproved},
		"labelsFirst":     githubv4.Int(100),
	}

	var response []*PullRequest
	for {
		if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
			return nil, err
		}
		for _, p := range query.Repository.PullRequests.Edges {
			labels := make([]LabelObject, len(p.Node.Labels.Edges))
			for _, l := range p.Node.Labels.Edges {
				labels = append(labels, l.Node.LabelObject)
			}

			for _, c := range p.Node.Commits.Edges {
				response = append(response, &PullRequest{
					PullRequestObject:   p.Node.PullRequestObject,
					Tip:                 c.Node.Commit,
					ApprovedReviewCount: p.Node.Reviews.TotalCount,
					Labels:              labels,
				})
			}
		}
		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		vars["prCursor"] = query.Repository.PullRequests.PageInfo.EndCursor
	}
	return response, nil
}

// GetPullRequest ...
func (m *GithubClient) GetPullRequest(prNumber int, commitRef string) (*PullRequest, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				PullRequestObject
				Commits struct {
					Edges []struct {
						Node struct {
							Commit CommitObject
						}
					}
				} `graphql:"commits(last:$commitsLast)"`
				Labels struct {
					Edges []struct {
						Node struct {
							LabelObject
						}
					}
				} `graphql:"labels(first:$labelsFirst)"`
			} `graphql:"pullRequest(number:$prNumber)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	vars := map[string]interface{}{
		"repositoryOwner": githubv4.String(m.Owner),
		"repositoryName":  githubv4.String(m.Repository),
		"prNumber":        githubv4.Int(prNumber),
		"commitsLast":     githubv4.Int(100),
		"labelsFirst":     githubv4.Int(100),
	}

	// TODO: Pagination - in case someone pushes > 100 commits before the build has time to start :p
	if err := m.V4.Query(context.TODO(), &query, vars); err != nil {
		return nil, err
	}

	labels := make([]LabelObject, len(query.Repository.PullRequest.Labels.Edges))
	for _, l := range query.Repository.PullRequest.Labels.Edges {
		labels = append(labels, l.Node.LabelObject)
	}
	for _, c := range query.Repository.PullRequest.Commits.Edges {
		if c.Node.Commit.OID == commitRef {
			// Return as soon as we find the correct ref.
			return &PullRequest{
				PullRequestObject: query.Repository.PullRequest.PullRequestObject,
				Tip:               c.Node.Commit,
				Labels:            labels,
			}, nil
		}
	}

	// Return an error if the commit was not found
	return nil, fmt.Errorf("commit with ref '%s' does not exist", commitRef)
}

// ListModifiedFiles in a pull request (not supported by V4 API).
func (m *GithubClient) ListModifiedFiles(prNumber int) ([]string, error) {
	var files []string

	opt := &github.ListOptions{
		PerPage: 100,
	}
	for {
		result, response, err := m.V3.PullRequests.ListFiles(
			context.TODO(),
			m.Owner,
			m.Repository,
			prNumber,
			opt,
		)
		if err != nil {
			return nil, err
		}
		for _, f := range result {
			files = append(files, *f.Filename)
		}
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	return files, nil
}

// PostComment to a pull request or issue.
func (m *GithubClient) PostComment(prNumber int, comment string) error {
	_, _, err := m.V3.Issues.CreateComment(
		context.TODO(),
		m.Owner,
		m.Repository,
		prNumber,
		&github.IssueComment{
			Body: github.String(comment),
		},
	)
	return err
}

// UpdateCommitStatus for a given commit (not supported by V4 API).
func (m *GithubClient) UpdateCommitStatus(commitRef, baseContext, statusContext, status, targetURL, description string) error {
	if baseContext == "" {
		baseContext = "concourse-ci"
	}

	if statusContext == "" {
		statusContext = "status"
	}

	if targetURL == "" {
		targetURL = strings.Join([]string{os.Getenv("ATC_EXTERNAL_URL"), "builds", os.Getenv("BUILD_ID")}, "/")
	}

	if description == "" {
		description = fmt.Sprintf("Concourse CI build %s", status)
	}

	_, _, err := m.V3.Repositories.CreateStatus(
		context.TODO(),
		m.Owner,
		m.Repository,
		commitRef,
		&github.RepoStatus{
			State:       github.String(strings.ToLower(status)),
			TargetURL:   github.String(targetURL),
			Description: github.String(description),
			Context:     github.String(path.Join(baseContext, statusContext)),
		},
	)
	return err
}

func (m *GithubClient) DeletePreviousComments(prNumber int) error {
	var getComments struct {
		Viewer struct {
			Login string
		}
		Repository struct {
			PullRequest struct {
				Id       string
				Comments struct {
					Edges []struct {
						Node struct {
							DatabaseId int64
							Author     struct {
								Login string
							}
						}
					}
				} `graphql:"comments(last:$commentsLast)"`
			} `graphql:"pullRequest(number:$prNumber)"`
		} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
	}

	vars := map[string]interface{}{
		"repositoryOwner": githubv4.String(m.Owner),
		"repositoryName":  githubv4.String(m.Repository),
		"prNumber":        githubv4.Int(prNumber),
		"commentsLast":    githubv4.Int(100),
	}

	if err := m.V4.Query(context.TODO(), &getComments, vars); err != nil {
		return err
	}

	for _, e := range getComments.Repository.PullRequest.Comments.Edges {
		if e.Node.Author.Login == getComments.Viewer.Login {
			_, err := m.V3.Issues.DeleteComment(context.TODO(), m.Owner, m.Repository, e.Node.DatabaseId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseRepository(s string) (string, string, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", errors.New("malformed repository")
	}
	return parts[0], parts[1], nil
}
