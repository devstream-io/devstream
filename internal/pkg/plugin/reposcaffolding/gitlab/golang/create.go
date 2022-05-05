package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create installs gitlab-repo-scaffolding-golang with provided options.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	return install(&opts)
}

func install(opts *rs.Options) (map[string]interface{}, error) {
	defer func() {
		if err := os.RemoveAll(DefaultWorkPath); err != nil {
			log.Errorf("Failed to clear workpath %s: %s.", DefaultWorkPath, err)
		}
	}()

	err := rs.CreateAndRenderLocalRepo(DefaultWorkPath, opts)
	if err != nil {
		return nil, err
	}

	if err := pushToRemote(filepath.Join(DefaultWorkPath, opts.Repo), opts); err != nil {
		return nil, err
	}

	return buildState(opts), nil
}
