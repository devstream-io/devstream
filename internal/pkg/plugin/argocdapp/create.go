package argocdapp

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates an ArgoCD app YAML and applys it.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	// validate parameters
	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// render an ArgoCD App YAML file based on inputs and template
	if err = writeContentToTmpFile(argoCDAppYAMLFile, argoCDAppTemplate, &opts); err != nil {
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
	return buildState(&opts), nil
}
