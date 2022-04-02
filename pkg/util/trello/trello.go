package trello

import (
	"fmt"
	"os"

	"github.com/adlio/trello"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type Client struct {
	*trello.Client
}

const DefaultListsNumber = 3

func NewClient() (*Client, error) {
	helpUrl := "https://docs.servicenow.com/bundle/quebec-it-asset-management/page/product/software-asset-management2/task/generate-trello-apikey-token.html"
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		return nil, fmt.Errorf("TRELLO_API_KEY and/or TRELLO_TOKEN are/is empty. see %s for more info", helpUrl)
	}

	return &Client{
		Client: trello.NewClient(apiKey, token),
	}, nil
}

func (c *Client) CreateBoard(kanbanBoardName, owner, repo string) (*trello.Board, error) {
	if kanbanBoardName == "" {
		kanbanBoardName = fmt.Sprintf("%s/%s", owner, repo)
	}

	board := trello.NewBoard(kanbanBoardName)
	board.Desc = boardDesc(owner, repo)

	err := c.Client.CreateBoard(&board, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (c *Client) CreateList(board *trello.Board, listName string) (*trello.List, error) {
	if listName == "" {
		return nil, fmt.Errorf("listName name can't be empty")
	}
	return c.Client.CreateList(board, listName, trello.Defaults())
}

// GetBoardIdAndListId get the board, which board name == kanbanBoardName, and board desc == owner/repo
func (c *Client) GetBoardIdAndListId(owner, repo, kanbanBoardName string) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	bs, err := c.Client.GetMyBoards()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		if checkTargetBoard(owner, repo, kanbanBoardName, b) {
			lists, err := b.GetLists()
			if err != nil {
				return nil, err
			}
			if len(lists) != DefaultListsNumber {
				log.Errorf("Unknown lists format: len==%d.", len(lists))
				return nil, fmt.Errorf("unknown lists format: len==%d", len(lists))
			}
			res["boardId"] = b.ID
			res["todoListId"] = lists[0].ID
			res["doingListId"] = lists[1].ID
			res["doneListId"] = lists[2].ID
		}
	}
	return res, nil
}

// CheckBoardExists check if board exists, which board name == kanbanBoardName, and board desc == owner/repo
func (c *Client) CheckBoardExists(owner, repo, kanbanBoardName string) (bool, error) {
	bs, err := c.Client.GetMyBoards()
	if err != nil {
		return false, err
	}

	for _, b := range bs {
		if checkTargetBoard(owner, repo, kanbanBoardName, b) {
			return true, nil
		}
	}
	return false, nil
}

// CheckAndDeleteBoard if the board exists, delete it
func (c *Client) CheckAndDeleteBoard(owner, repo, kanbanBoardName string) error {
	bs, err := c.Client.GetMyBoards()
	if err != nil {
		return err
	}

	for _, b := range bs {
		if checkTargetBoard(owner, repo, kanbanBoardName, b) {
			log.Infof("Board will be deleted, name: %s, description: %s.", b.Name, b.Desc)
			return b.Delete()
		}
	}
	return nil
}

func checkTargetBoard(owner, repo, kanbanBoardName string, b *trello.Board) bool {
	return !b.Closed && b.Name == kanbanBoardName && b.Desc == boardDesc(owner, repo)
}

func boardDesc(owner, repo string) string {
	return fmt.Sprintf("Description is managed by DevStream, please don't modify. %s/%s", owner, repo)
}
