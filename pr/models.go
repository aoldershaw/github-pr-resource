package pr

import (
	"errors"

	"github.com/cloudfoundry-community/github-pr-instances-resource/models"
)

// Source represents the configuration for the resource.
type Source struct {
	models.CommonConfig
	models.GithubConfig
	Number        int      `json:"number"`
	GitCryptKey   string   `json:"git_crypt_key"`
	DisableGitLFS bool     `json:"disable_git_lfs"`
	Paths         []string `json:"paths"`
	IgnorePaths   []string `json:"ignore_paths"`
	DisableCISkip bool     `json:"disable_ci_skip"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AccessToken == "" {
		return errors.New("access_token must be set")
	}
	if s.Repository == "" {
		return errors.New("repository must be set")
	}

	isHostingEndpointEnabled := s.HostingEndpoint != ""
	isV3EndpointEnabled := s.V3Endpoint != ""
	isV4EndpointEnabled := s.V4Endpoint != ""
	if isHostingEndpointEnabled != isV3EndpointEnabled || isV3EndpointEnabled != isV4EndpointEnabled {
		return errors.New("if any of hosting_endpoint, v3_endpoint, or v4_endpoint are set, all of them must be set")
	}
	return nil
}

type Version struct {
	Ref string `json:"ref"`
}
