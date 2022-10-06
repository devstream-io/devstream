package mapz

import (
	"github.com/mitchellh/mapstructure"
)

func FillMapWithStrAndError(keys []string, value error) map[string]error {
	retMap := make(map[string]error, len(keys))
	if len(keys) == 0 {
		return retMap
	}

	for _, key := range keys {
		retMap[key] = value
	}
	return retMap
}

func DecodeStructToMap(structVars any) (map[string]interface{}, error) {
	var rawConfigVars map[string]interface{}
	if err := mapstructure.Decode(structVars, &rawConfigVars); err != nil {
		return nil, err
	}
	return rawConfigVars, nil
}
