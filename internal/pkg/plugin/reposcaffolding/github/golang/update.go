package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
)

// Update re-installs github-repo-scaffolding-golang with provided options.
func Update(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	_, err := uninstall(&param)
	if err != nil {
		return nil, err
	}

	return install(&param)
}
