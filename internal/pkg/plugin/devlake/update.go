package devlake

func Update(options *map[string]interface{}) (map[string]interface{}, error) {
	return Create(options)
}
