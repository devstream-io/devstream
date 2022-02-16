package argocdapp

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/kubectl"
)

// Create creates an ArgoCD app YAML and applys it.
func Create(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param

	// decode input parameters into a struct
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return nil, err
	}

	// validate parameters
	if errs := validateParams(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	// render an ArgoCD App YAML file based on inputs and template
	if err = writeContentToTmpFile(argoCDAppYAMLFile, argoCDAppTemplate, &param); err != nil {
		return nil, err
	}

	// kubectl apply -f
	if err = kubectl.KubeApply(argoCDAppYAMLFile); err != nil {
		return nil, err
	}

	// remove temporary YAML file used for kubectl apply
	if err = os.Remove(argoCDAppYAMLFile); err != nil {
		log.Warnf("Temporary YAML file %s can't be deleted, but the installation is successful.", argoCDAppYAMLFile)
	}

	// build state & return results
	return buildState(param), nil
}
