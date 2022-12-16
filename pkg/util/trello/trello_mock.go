package trello

import "github.com/adlio/trello"

type MockTrelloClient struct {
	GetError    error
	GetValue    *trello.Board
	CreateError error
}

func (c *MockTrelloClient) Create(name, description string) (*trello.Board, error) {
	return nil, c.CreateError
}

func (c *MockTrelloClient) Get(name, description string) (*trello.Board, error) {
	return c.GetValue, c.GetError
}
