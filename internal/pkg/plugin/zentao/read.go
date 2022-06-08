package zentao

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
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

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	ready, err := CheckDeploymentsAndServicesReady(kubeClient, &opts)
	if err != nil {
		return nil, err
	}

	if !ready {
		return map[string]interface{}{"stopped": true}, nil
	}

	return map[string]interface{}{"running": true}, nil
}
