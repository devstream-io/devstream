package trello

func Read(options map[string]interface{}) (map[string]interface{}, error) {

	var opt *Options
	var err error

	if opt, err = convertMap2Options(options); err != nil {
		return nil, err
	}

	if err := validateOptions(opt); err != nil {
		return nil, err
	}

	return buildReadState(opt)
}
