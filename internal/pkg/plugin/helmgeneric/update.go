package helmgeneric

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// 2. install or upgrade
	if err := InstallOrUpgradeChart(&opts); err != nil {
		return nil, err
	}

	// 3. fill the return map
	retMap := make(map[string]interface{})
	log.Debugf("Return map: %v", retMap)

	return nil, nil
}
