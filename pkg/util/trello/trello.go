package trello

import (
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/merico-dev/stream/internal/pkg/log"
)

type Client struct {
	*trello.Client
}

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

func (c *Client) CreateBoard(boardName string) (*trello.Board, error) {
	if boardName == "" {
		return nil, fmt.Errorf("board name can't be empty")
	}
	board := trello.NewBoard(boardName)
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

func (c *Client) GetBoardIdAndListId() (map[string]interface{}, error) {
	res := make(map[string]interface{})

	bs, err := c.Client.GetMyBoards()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		lists, err := b.GetLists()
		if err != nil {
			return nil, err
		}
		if len(lists) != 3 {
			log.Errorf("Unknown lists format: len==%d", len(lists))
			return nil, fmt.Errorf("unknown lists format: len==%d", len(lists))
		}
		res["boardId"] = b.ID
		res["todoListId"] = lists[0].ID
		res["doingListId"] = lists[1].ID
		res["doneListId"] = lists[2].ID
	}
	return res, nil
}
