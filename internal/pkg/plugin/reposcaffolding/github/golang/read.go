package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/github"
)

// Read check the health for github-repo-scaffolding-golang with provided options.
func Read(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	return check(&param)
}

func check(param *Param) (map[string]interface{}, error) {
	ghOptions := &github.Option{
		Owner:    param.Owner,
		Repo:     param.Repo,
		NeedAuth: true,
	}

	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	return ghClient.GetRepoInfo()
}
