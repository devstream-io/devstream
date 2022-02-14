package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
)

// Reinstall re-installs github-repo-scaffolding-golang with provided options.
func Reinstall(options *map[string]interface{}) (bool, error) {
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

	_, err := uninstall(&param)
	if err != nil {
		return false, err
	}

	_, err = install(&param)
	if err != nil {
		return false, err
	}
	log.Successf("GitHub repo %s reinstalled.", param.Repo)

	return true, nil
}
