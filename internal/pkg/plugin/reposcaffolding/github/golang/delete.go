package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Delete uninstalls github-repo-scaffolding-golang with provided options.
func Delete(param map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(param, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("options are illegal")
	}

	return uninstall(&opts)
}

func uninstall(opts *Options) (bool, error) {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Repo:     opts.Repo,
		NeedAuth: true,
	}

	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}
	if err := ghClient.DeleteRepo(); err != nil {
		return false, err
	}
	return true, nil
}
