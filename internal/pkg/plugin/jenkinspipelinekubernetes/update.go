package jenkinspipelinekubernetes

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := ValidateAndDefaults(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// TODO(aFlyBird0): determine how to update the resource, such as:
	// if some config/resource are changed, we should restart the Jenkins
	// some, we should only call some update function
	// others, we just ignore them

	// now we just use the same way as create,
	// because the logic is the same: "if not exists, create; if exists, do nothing"
	// if it changes in the future, we should change the way to update
	return Create(options)
}
