package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update re-installs github-repo-scaffolding-golang with provided options.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	var opts rs.Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := rs.Validate(&opts); len(errs) != 0 {
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
