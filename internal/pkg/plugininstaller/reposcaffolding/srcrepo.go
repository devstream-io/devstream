package reposcaffolding

import (
	"fmt"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/zip"
)

// default get main branch of repo for scaffolding project
var srcDefaultBranch = "main"

// SrcRepo describe how to get scaffolding repo
type SrcRepo struct {
	Repo     string `validate:"required"`
	Org      string `validate:"required"`
	RepoType string `validate:"oneof=gitlab github" mapstructure:"repo_type"`
}

func (t *SrcRepo) DownloadRepo(workpath string) error {
	// 1. download scaffolding repo from github
	if err := downloadGithubRepo(t.Org, t.Repo, workpath); err != nil {
		return err
	}

	// 2. unzip downloaded zip file
	unzipPath := filepath.Join(workpath, github.DefaultLatestCodeZipfileName)
	if err := zip.UnZip(unzipPath, workpath); err != nil {
		return err
	}
	return nil
}

func (t *SrcRepo) getLocalRepoPath(workpath string) string {
	return filepath.Join(workpath, fmt.Sprintf("%s-%s", t.Repo, srcDefaultBranch))
}
