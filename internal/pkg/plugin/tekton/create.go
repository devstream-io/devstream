package tekton

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	if err := DealWithNsWhenInstall(&opts); err != nil {
		return nil, err
	}
	var retErr error
	// delete namespace if encounter error
	defer func() {
		if retErr == nil {
			return
		}
		if err := DealWithNsWhenInterruption(&opts); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")
	}()

	if retErr = InstallOrUpgradeChart(&opts); retErr != nil {
		return nil, retErr
	}

	retMap := GetStaticState().ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)
	return retMap, nil
}
