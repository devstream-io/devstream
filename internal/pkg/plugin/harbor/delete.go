package harbor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	// 1. create helm instance from params
	h, err := helm.NewHelm(opts.GetHelmParam())
	if err != nil {
		return false, err
	}
	log.Info("Uninstalling harbor helm chart.")

	// 2. delete harbor by helm
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	// 3. delete ns if helm is deleted
	if err := DealWithNsWhenInterruption(&opts); err != nil {
		return false, err
	}
	return true, nil

}
