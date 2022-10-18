package trellogithub

import "github.com/devstream-io/devstream/internal/pkg/configmanager"

// Delete remove trello-github-integ workflows.
func Delete(options configmanager.RawOptions) (bool, error) {
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
