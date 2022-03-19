package openldap

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/merico-dev/stream/internal/pkg/plugin/common/helm"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create creates OpenLDAP with provided options.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. decode options
	var param Param
	if err := mapstructure.Decode(options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	// 2. deal with ns
	if err := DealWithNsWhenInstall(&param); err != nil {
		return nil, err
	}

	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := DealWithNsWhenInterruption(&param); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")
	}()

	// 3. install or upgrade
	if retErr = InstallOrUpgradeChart(&param); retErr != nil {
		return nil, retErr
	}

	// 4. fill the return map
	retMap := GetStaticState().ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}
