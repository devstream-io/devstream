package trellogithub

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
)

// Delete remove trello-github-integ workflows.
func Delete(options configmanager.RawOptions) (bool, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

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
