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
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var param Param
	if err := mapstructure.Decode(options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	return install(&param)
}

func install(param *Param) (map[string]interface{}, error) {
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

	if err := push(param); err != nil {
		return nil, err
	}

	return buildState(param), nil
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

func push(param *Param) error {
	ghOption := &github.Option{
		Owner:    param.Owner,
		Repo:     param.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return err
	}

	err = InitRepoLocalAndPushToRemote(filepath.Join(DefaultWorkPath, DefaultTemplateRepo+"-main"), param, ghClient)
	if err != nil {
		return err
	}

	return nil
}
