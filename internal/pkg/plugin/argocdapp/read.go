package argocdapp

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
)

func Read(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param

	// decode input parameters into a struct
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return nil, err
	}

	// validate parameters
	if errs := validateParams(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	// describe app with retry
	state := make(map[string]interface{})
	operation := func() error {
		err := getArgoCDAppFromK8sAndSetState(state, param.App.Name, param.App.Namespace)
		if err != nil {
			return err
		}
		return nil
	}
	bkoff := backoff.NewExponentialBackOff()
	bkoff.MaxElapsedTime = 3 * time.Minute
	err = backoff.Retry(operation, bkoff)
	if err != nil {
		return nil, err
	}

	return state, nil
}
