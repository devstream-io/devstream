package trello

func Read(options map[string]interface{}) (map[string]interface{}, error) {

	var opts *Options
	var err error

	if opts, err = convertMap2Options(options); err != nil {
		return nil, err
	}

	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	return buildReadState(opts)
}
