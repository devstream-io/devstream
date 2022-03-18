package trellogithub

// Delete remove trello-github-integ workflows.
func Delete(options map[string]interface{}) (bool, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return false, err
	}

	err = tg.client.DeleteWorkflow(trelloWorkflow, tg.options.Branch)
	if err != nil {
		return false, err
	}

	return true, nil
}
