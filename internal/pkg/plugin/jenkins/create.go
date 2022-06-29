package jenkins

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates jenkins with provided options.
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
	if err := helm.DealWithNsWhenInstall(&opts.Options); err != nil {
		return nil, err
	}

	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := helm.DealWithNsWhenInterruption(&opts.Options); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")

		// Clear all the resources have been created if the creation process interruption.
		if err := postDelete(); err != nil {
			log.Errorf("Failed to clear the resources have been created: %s.", err)
		}
	}()

	// 3. pre-create
	if retErr = preCreate(opts); retErr != nil {
		log.Errorf("The pre-create logic failed: %s.", retErr)
		return nil, retErr
	}

	// 4. install or upgrade
	if retErr = helm.InstallOrUpgradeChart(&opts.Options); retErr != nil {
		return nil, retErr
	}

	// 5. fill the return map
	releaseName := opts.Chart.ReleaseName
	retMap := GetStaticState(releaseName).ToStringInterfaceMap()
	log.Debugf("Return map: %v.", retMap)

	// show how to get pwd of the admin user
	howToGetPasswdOfAdmin(&opts)

	// show jenkins url
	showJenkinsUrl(&opts)

	return retMap, nil
}
