package hashicorpvault

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates hashicorp-vault with provided options.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. decode options
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

	// 2. deal with ns
	if err := DealWithNsWhenInstall(&opts); err != nil {
		return nil, err
	}

	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := DealWithNsWhenInterruption(&opts); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")
	}()

	// 3. install or upgrade
	if retErr = InstallOrUpgradeChart(&opts); retErr != nil {
		return nil, retErr
	}

	// 4. fill the return map
	retMap := GetStaticState().ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}
