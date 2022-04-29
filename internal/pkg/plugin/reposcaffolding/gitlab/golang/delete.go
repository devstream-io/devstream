package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

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
	c, err := gitlab.NewClient()
	if err != nil {
		return false, err
	}

	if err := c.DeleteProject(opts.PathWithNamespace); err != nil {
		log.Errorf("Failed to create repo: %s.", err)
		return false, err
	}

	return true, nil
}
