package trello

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/spf13/viper"
)

type Client struct {
	*trello.Client
}

func NewClient() (*Client, error) {
	helpUrl := "https://docs.servicenow.com/bundle/quebec-it-asset-management/page/product/software-asset-management2/task/generate-trello-apikey-token.html"
	apiKey := viper.GetString("TRELLO_API_KEY")
	token := viper.GetString("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		return nil, fmt.Errorf("TRELLO_API_KEY and/or TRELLO_TOKEN are/is empty. see %s for more info", helpUrl)
	}

	return &Client{
		Client: trello.NewClient(apiKey, token),
	}, nil
}

func (c *Client) CreateBoard(boardName string) error {
	if boardName == "" {
		return fmt.Errorf("board name can't be empty")
	}
	board := trello.NewBoard(boardName)
	return c.Client.CreateBoard(&board, trello.Defaults())
}
