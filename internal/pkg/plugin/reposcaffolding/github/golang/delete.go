package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Delete uninstalls github-repo-scaffolding-golang with provided options.
func Delete(options map[string]interface{}) (bool, error) {
	var opts rs.Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := rs.Validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	return uninstall(&opts)
}

func uninstall(opts *rs.Options) (bool, error) {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
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
