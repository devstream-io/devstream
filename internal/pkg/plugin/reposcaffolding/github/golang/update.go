package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
)

// Update re-installs github-repo-scaffolding-golang with provided options.
func Update(params map[string]interface{}) (map[string]interface{}, error) {
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

	_, err := uninstall(&opts)
	if err != nil {
		return nil, err
	}

	return install(&opts)
}
