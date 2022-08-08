package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/github"
)

// default get main branch of repo for scaffolding project
var srcDefaultBranch = "main"

// SrcRepo describe how to get scaffolding repo
type SrcRepo struct {
	Repo     string `validate:"required"`
	Org      string `validate:"required"`
	RepoType string `validate:"oneof=github" mapstructure:"repo_type"`
}

func (t *SrcRepo) getRepoName() string {
	return fmt.Sprintf("%s-%s", t.Repo, srcDefaultBranch)
}

func (t *SrcRepo) getDownloadURL() (string, error) {
	ghOption := &github.Option{
		Org:      t.Org,
		Repo:     t.Repo,
		NeedAuth: false,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return "", err
	}
	return ghClient.GetLatestCodeZipURL(), nil
}
