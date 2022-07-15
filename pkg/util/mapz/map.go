package mapz

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
