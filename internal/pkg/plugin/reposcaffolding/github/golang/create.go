package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/zip"
)

// Create installs github-repo-scaffolding-golang with provided options.
func Create(params map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(params, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("options are illegal")
	}

	return install(&opts)
}

func install(opts *Options) (map[string]interface{}, error) {
	// Clear workpath before return
	defer func() {
		if err := os.RemoveAll(DefaultWorkPath); err != nil {
			log.Errorf("Failed to clear workpath %s: %s.", DefaultWorkPath, err)
		}
	}()

	if err := download(); err != nil {
		return nil, err
	}

	if err := zip.UnZip(filepath.Join(DefaultWorkPath, github.DefaultLatestCodeZipfileName), DefaultWorkPath); err != nil {
		return nil, err
	}

	if err := push(opts); err != nil {
		return nil, err
	}

	return buildState(opts), nil
}

func download() error {
	ghOption := &github.Option{
		Owner:    DefaultTemplateOwner,
		Repo:     DefaultTemplateRepo,
		NeedAuth: false,
		WorkPath: DefaultWorkPath,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return err
	}

	if err = ghClient.DownloadLatestCodeAsZipFile(); err != nil {
		return err
	}

	return nil
}

func push(opts *Options) error {
	ghOption := &github.Option{
		Owner:    opts.Owner,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return err
	}

	err = InitRepoLocalAndPushToRemote(filepath.Join(DefaultWorkPath, DefaultTemplateRepo+"-main"), opts, ghClient)
	if err != nil {
		return err
	}

	return nil
}
