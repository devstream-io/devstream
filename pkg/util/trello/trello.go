package trello

import (
	"fmt"
	"os"

	"github.com/adlio/trello"
)

type TrelloAPI interface {
	Create(name, description string) (*trello.Board, error)
	Get(name, description string) (*trello.Board, error)
}

type client struct {
	*trello.Client
}

func NewClient() (TrelloAPI, error) {
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		const helpUrl = "https://docs.servicenow.com/bundle/quebec-it-asset-management/page/product/software-asset-management2/task/generate-trello-apikey-token.html"
		return nil, fmt.Errorf("TRELLO_API_KEY and/or TRELLO_TOKEN are/is empty. see %s for more info", helpUrl)
	}

	return &client{
		Client: trello.NewClient(apiKey, token),
	}, nil
}

func (c *client) Create(name, description string) (*trello.Board, error) {
	board := trello.NewBoard(name)
	board.Desc = description

	err := c.Client.CreateBoard(&board, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (c *client) Get(name, description string) (*trello.Board, error) {
	bs, err := c.Client.GetMyBoards()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		if !b.Closed && b.Name == name && b.Desc == description {
			return b, nil
		}
	}
	return nil, nil
}
