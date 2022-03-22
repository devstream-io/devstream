package trello

// Delete delete trello board and lists
func Delete(options map[string]interface{}) (bool, error) {
	var opts *Options
	var err error

	if opts, err = convertMap2Options(options); err != nil {
		return false, err
	}
	if err := validateOptions(opts); err != nil {
		return false, err
	}

	if err = DeleteTrelloBoard(opts); err != nil {
		return false, err
	}
	return true, nil
}
