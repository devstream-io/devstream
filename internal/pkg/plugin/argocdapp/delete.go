package argocdapp

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/kubectl"
	"github.com/merico-dev/stream/pkg/util/log"
)

func Delete(options *map[string]interface{}) (bool, error) {
	var param Param

	// decode input parameters into a struct
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	// validate parameters
	if errs := validateParams(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	// render an ArgoCD App YAML file based on inputs and template
	if err = writeContentToTmpFile(argoCDAppYAMLFile, argoCDAppTemplate, &param); err != nil {
		return false, err
	}

	// kubectl delete -f
	if err = kubectl.KubeDelete(argoCDAppYAMLFile); err != nil {
		return false, err
	}

	// remove temporary YAML file used for kubectl apply
	if err = os.Remove(argoCDAppYAMLFile); err != nil {
		log.Warnf("Temporary YAML file %s can't be deleted, but the installation is successful.", argoCDAppYAMLFile)
	}

	return true, nil
}
