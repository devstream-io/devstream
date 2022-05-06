package helmgeneric

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts helm.Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	if err := helm.DealWithNsWhenInstall(&opts); err != nil {
		return nil, err
	}

	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := helm.DealWithNsWhenInterruption(&opts); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")
	}()

	// install or upgrade
	if retErr = helm.InstallOrUpgradeChart(&opts); retErr != nil {
		return nil, retErr
	}

	// fill the return map
	retMap := make(map[string]interface{})
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}
