package devlake

import (
	"github.com/mitchellh/mapstructure"
)

const (
	DevLakeTotalK8sDeployments = 4
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	return readDeploymentsAndServicesAndBuildState()
}
