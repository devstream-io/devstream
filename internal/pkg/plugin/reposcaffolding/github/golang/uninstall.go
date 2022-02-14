package golang

import (
	"fmt"
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/mitchellh/mapstructure"
)

// Uninstall uninstalls github-repo-scaffolding-golang with provided options.
func Uninstall(options *map[string]interface{}) (bool, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	return uninstall(&param)
}

func uninstall(param *Param) (bool, error) {
	ghOptions := &github.Option{
		Owner:    param.Owner,
		Repo:     param.Repo,
		NeedAuth: true,
	}

	gitHubClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}
	err = gitHubClient.DeleteRepo()
	if err != nil {
		return false, err
	}
	return true, nil
}
