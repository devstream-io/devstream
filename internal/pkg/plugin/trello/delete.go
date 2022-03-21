package trello

// Delete delete trello board and lists
func Delete(options map[string]interface{}) (bool, error) {
	var opt *Options
	var err error

	if opt, err = convertMap2Options(options); err != nil {
		return false, err
	}
	if err = validateOptions(opt); err != nil {
		return false, err
	}

	if err = DeleteTrelloBoard(opt); err != nil {
		return false, err
	}
	return true, nil
}
