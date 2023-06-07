package test_helpers

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/cloudfoundry-community/github-pr-instances-resource/models"
	"github.com/shurcooL/githubv4"
)

func CreateTestPR(
	count int,
	baseName string,
	skipCI bool,
	isCrossRepo bool,
	approvedReviews int,
	labels []string,
	isDraft bool,
	state githubv4.PullRequestState,
) *models.PullRequest {
	n := strconv.Itoa(count)
	d := time.Now().AddDate(0, 0, -count)
	m := fmt.Sprintf("commit message%s", n)
	if skipCI {
		m = "[skip ci]" + m
	}
	approvedCount := approvedReviews

	var labelObjects []models.LabelObject
	for _, l := range labels {
		lObject := models.LabelObject{
			Name: l,
		}

		labelObjects = append(labelObjects, lObject)
	}

	return &models.PullRequest{
		PullRequestObject: models.PullRequestObject{
			ID:          fmt.Sprintf("pr%s", n),
			Number:      count,
			Title:       fmt.Sprintf("pr%s title", n),
			URL:         fmt.Sprintf("pr%s url", n),
			BaseRefName: baseName,
			HeadRefName: fmt.Sprintf("pr%s", n),
			Repository: struct{ URL string }{
				URL: fmt.Sprintf("repo%s url", n),
			},
			IsCrossRepository: isCrossRepo,
			IsDraft:           isDraft,
			State:             state,
			ClosedAt:          githubv4.DateTime{Time: time.Now()},
			MergedAt:          githubv4.DateTime{Time: time.Now()},
		},
		Tip: models.CommitObject{
			ID:            fmt.Sprintf("commit%s", n),
			OID:           fmt.Sprintf("oid%s", n),
			CommittedDate: githubv4.DateTime{Time: d},
			Message:       m,
			Author: struct {
				User  struct{ Login string }
				Email string
			}{
				User: struct{ Login string }{
					Login: fmt.Sprintf("login%s", n),
				},
				Email: "user@example.com",
			},
		},
		ApprovedReviewCount: approvedCount,
		Labels:              labelObjects,
	}
}

func CreateTestDirectory(t *testing.T) string {
	dir, err := ioutil.TempDir("", "github-pr-resource")
	if err != nil {
		t.Fatalf("failed to create temporary directory")
	}
	return dir
}

func ReadTestFile(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read: %s: %s", path, err)
	}
	return string(b)
}
