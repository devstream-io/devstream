package github

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	RegistryDockerhub   RegistryType = "dockerhub"
	RegistryHarbor      RegistryType = "harbor"
	CommitMessage                    = "GitHub Actions workflow, created by DevStream"
	PRBuilderFileName                = "pr-builder.yml"
	MainBuilderFileName              = "main-builder.yml"
)

// GithubActionLanguage is the language of repo
type GithubActionLanguage struct {
	Name    string `validate:"required"`
	Version string
}

// GithubActionOptions is the struct for configurations of the  plugin.
type GithubActionOptions struct {
	Owner    string                `validate:"required_without=Org"`
	Org      string                `validate:"required_without=Owner"`
	Repo     string                `validate:"required"`
	Branch   string                `validate:"required"`
	Language *GithubActionLanguage `validate:"required"`

	// optional
	Workflows []*github.Workflow
	Build     Build
	Test      *Test
	Docker    *Docker
}

func NewGithubActionOptions(options configmanager.RawOptions) (*GithubActionOptions, error) {
	var opts GithubActionOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (opts *GithubActionOptions) GetLanguage() string {
	return fmt.Sprintf("%s-%s", opts.Language.Name, opts.Language.Version)
}

func (opts *GithubActionOptions) GetGithubClient() (*github.Client, error) {
	ghOptions := &git.RepoInfo{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	return github.NewClient(ghOptions)
}

func (opts *GithubActionOptions) RenderWorkFlow(content string) (string, error) {
	return template.Render("github-actions", content, opts)
}

func (opts *GithubActionOptions) CheckAddDockerHubToken() bool {
	if opts.Docker != nil && opts.Docker.Registry.Type == "dockerhub" {
		return true
	}
	return false
}
