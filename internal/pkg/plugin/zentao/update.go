package zentao

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
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

	if err := updateByClientAPI(&opts); err != nil {
		return nil, err
	}

	return map[string]interface{}{"running": true}, nil
}

func updateByClientAPI(opts *Options) error {
	log.Debug("Start zentao plugin update.")
	// 1. Create k8s clientset
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// 2. Delete zentao application
	if err = DeleteZentaoAPP(kubeClient, opts); err != nil {
		return err
	}

	// 3. Recreate zentao application
	if err = CreateZentaoAPP(kubeClient, opts); err != nil {
		return err
	}

	return nil
}
