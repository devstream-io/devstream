package reposcaffolding

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/zip"
)

const (
	TemplateRepo       = "dtm-scaffolding-golang"
	TemplateOrg        = "devstream-io"
	TransitBranch      = "init-with-devstream"
	MainBranch         = "main"
	AppNamePlaceHolder = "_app_name_"
)

type Config struct {
	AppName   string
	ImageRepo string
	Repo      Repo
}

type Repo struct {
	Name  string
	Owner string
}

func CreateAndRenderLocalRepo(workpath string, opts *Options) error {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return err
	}

	if err := download(TemplateOrg, TemplateRepo, workpath); err != nil {
		return err
	}

	if err := zip.UnZip(filepath.Join(workpath, github.DefaultLatestCodeZipfileName), workpath); err != nil {
		return err
	}

	if retErr := walkLocalRepoPath(workpath, opts, ghClient); retErr != nil {
		log.Debugf("Failed to walk local repo-path: %s.", retErr)
		return retErr
	}

	return nil
}
