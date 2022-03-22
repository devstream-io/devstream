package trello

// Update recreate trello board and lists.
func Update(options map[string]interface{}) (map[string]interface{}, error) {

	var opts *Options
	var err error

	if opts, err = convertMap2Options(options); err != nil {
		return nil, err
	}

	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	if err = DeleteTrelloBoard(opts); err != nil {
		return nil, err
	}

	trelloIds, err := CreateTrelloBoard(opts)
	if err != nil {
		return nil, err
	}

	return buildState(opts, trelloIds), nil
}
