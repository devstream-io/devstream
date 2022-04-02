package argocdapp

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return false, err
	}

	// validate parameters
	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	// render an ArgoCD App YAML file based on inputs and template
	if err = writeContentToTmpFile(argoCDAppYAMLFile, argoCDAppTemplate, &opts); err != nil {
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
