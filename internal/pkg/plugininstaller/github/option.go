package github

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/template"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/github"
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

func NewGithubActionOptions(options plugininstaller.RawOptions) (*GithubActionOptions, error) {
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
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	return github.NewClient(ghOptions)
}

func (opts *GithubActionOptions) RenderWorkFlow(content string) (string, error) {
	return template.Render("githubactions", content, opts)
}

func (opts *GithubActionOptions) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func (opts *GithubActionOptions) CheckAddDockerHubToken() bool {
	if opts.Docker != nil && opts.Docker.Registry.Type == "dockerhub" {
		return true
	}
	return false
}
