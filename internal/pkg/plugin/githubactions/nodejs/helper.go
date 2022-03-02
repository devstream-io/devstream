package nodejs

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
)

func parseAndValidateOptions(options map[string]interface{}) (*Options, error) {
	var opt Options
	err := mapstructure.Decode(options, &opt)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opt); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("incorrect params")
	}

	return &opt, nil
}
