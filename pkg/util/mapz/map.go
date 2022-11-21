package mapz

import (
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/maps"
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

// Merge merge two maps
// if there are same keys in two maps, the key of second one will overwrite the first one
func Merge(m1 map[string]any, m2 map[string]any) map[string]any {
	m1Clone := maps.Clone(m1)
	maps.Copy(m1Clone, m2)
	return m1Clone
}
