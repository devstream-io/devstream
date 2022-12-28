package trello

import (
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/trello"
)

func createBoard(rawOptions configmanager.RawOptions) error {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return err
	}
	c, err := trello.NewClient()
	if err != nil {
		return err
	}
	return opts.Board.create(c)

}

// deleteBoard delete specified board
func deleteBoard(rawOptions configmanager.RawOptions) error {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return err
	}
	c, err := trello.NewClient()
	if err != nil {
		return err
	}
	return opts.Board.delete(c)
}

// addTrelloSecret will add trello secret in github
func addTrelloSecret(rawOptions configmanager.RawOptions) error {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return err
	}

	// 1. init scm client
	scmClient, err := scm.NewClientWithAuth(opts.Scm)
	if err != nil {
		return err
	}

	// 2. init trello client
	trelloClient, err := trello.NewClient()
	if err != nil {
		return err
	}
	trelloBoard, err := opts.Board.get(trelloClient)
	if err != nil {
		return err
	}

	// 3. add github secret
	// add key
	if err := scmClient.AddRepoSecret("TRELLO_API_KEY", viper.GetString("trello_api_key")); err != nil {
		return err
	}
	// add token
	if err := scmClient.AddRepoSecret("TRELLO_TOKEN", viper.GetString("trello_token")); err != nil {
		return err
	}
	// add board id
	if err := scmClient.AddRepoSecret("TRELLO_BOARD_ID", trelloBoard.boardID); err != nil {
		return err
	}
	// add todolist id
	if err := scmClient.AddRepoSecret("TRELLO_TODO_LIST_ID", trelloBoard.todoListID); err != nil {
		return err
	}
	// add doinglist id
	if err := scmClient.AddRepoSecret("TRELLO_DOING_LIST_ID", trelloBoard.doingListID); err != nil {
		return err
	}
	// add donelist id
	if err := scmClient.AddRepoSecret("TRELLO_DONE_LIST_ID", trelloBoard.doneListID); err != nil {
		return err
	}
	// add member map
	if err := scmClient.AddRepoSecret("TRELLO_MEMBER_MAP", "[]"); err != nil {
		return err
	}
	return nil
}
