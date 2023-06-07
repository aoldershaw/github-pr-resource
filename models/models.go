package models

import "github.com/shurcooL/githubv4"

// Metadata output from get/put steps.
type Metadata []*MetadataField

// Add a MetadataField to the Metadata.
func (m *Metadata) Add(name, value string) {
	*m = append(*m, &MetadataField{Name: name, Value: value})
}

// MetadataField ...
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PullRequest represents a pull request and includes the tip (commit).
type PullRequest struct {
	PullRequestObject
	Tip                 CommitObject
	ApprovedReviewCount int
	Labels              []LabelObject
}

// PullRequestObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/pullrequest/
type PullRequestObject struct {
	ID          string
	Number      int
	Title       string
	URL         string
	BaseRefName string
	HeadRefName string
	Repository  struct {
		URL string
	}
	IsCrossRepository bool
	IsDraft           bool
	State             githubv4.PullRequestState
	ClosedAt          githubv4.DateTime
	MergedAt          githubv4.DateTime
}

// UpdatedDate returns the last time a PR was updated, either by commit
// or being closed/merged.
func (p *PullRequest) UpdatedDate() githubv4.DateTime {
	date := p.Tip.CommittedDate
	switch p.State {
	case githubv4.PullRequestStateClosed:
		date = p.ClosedAt
	case githubv4.PullRequestStateMerged:
		date = p.MergedAt
	}
	return date
}

// CommitObject represents the GraphQL commit node.
// https://developer.github.com/v4/object/commit/
type CommitObject struct {
	ID            string
	OID           string
	CommittedDate githubv4.DateTime
	Message       string
	Author        struct {
		User struct {
			Login string
		}
		Email string
	}
}

// ChangedFileObject represents the GraphQL FilesChanged node.
// https://developer.github.com/v4/object/pullrequestchangedfile/
type ChangedFileObject struct {
	Path string
}

// LabelObject represents the GraphQL label node.
// https://developer.github.com/v4/object/label
type LabelObject struct {
	Name string
}
