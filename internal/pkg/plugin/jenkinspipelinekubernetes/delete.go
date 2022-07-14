package jenkinspipelinekubernetes

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

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

	// TODO(aFlyBird0): use helm uninstall to delete the resource created by devstream
	// TODO(aFlyBird0): filter ConfigMaps Created by devstream in the "jenkins" namespace and delete them(filter by label)
	// TODO(aFlyBird0): delete other resources created by devstream(if exists)

	return true, nil
}
