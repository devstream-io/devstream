package trellogithub

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/pkg/util/github"
)

// VerifyWorkflows get the workflows with names "wf1.yml", "wf2.yml", then:
// If all workflows is ok => return ({"wf1.yml":nil, "wf2.yml:nil}, nil)
// If some error occurred => return (nil, error)
// If wf1.yml is not found => return ({"wf1.yml":error("not found"), "wf2.yml:nil},nil)
func (tg *TrelloGithub) VerifyWorkflows(workflows []*github.Workflow) (map[string]error, error) {
	wsFiles := make([]string, 0)
	for _, w := range workflows {
		wsFiles = append(wsFiles, w.WorkflowFileName)
	}

	fmt.Printf("Workflow files: %v", wsFiles)
	filesInRemoteDir, rMap, err := tg.client.FetchRemoteContent(wsFiles)
	if err != nil {
		return nil, err
	}
	if rMap != nil {
		return rMap, nil
	}

	return tg.CompareFiles(wsFiles, filesInRemoteDir), nil
}

// AddTrelloIdSecret add trello ids to secret
func (tg *TrelloGithub) AddTrelloIdSecret(trelloId *TrelloItemId) error {
	ghOptions := &github.Option{
		Owner:    tg.options.Owner,
		Repo:     tg.options.Repo,
		NeedAuth: true,
	}
	c, err := github.NewClient(ghOptions)
	if err != nil {
		return err
	}
	// add key
	if err := c.AddRepoSecret("TRELLO_API_KEY", viper.GetString("trello_api_key")); err != nil {
		return err
	}

	// add token
	if err := c.AddRepoSecret("TRELLO_TOKEN", viper.GetString("trello_token")); err != nil {
		return err
	}

	// add board id
	if err := c.AddRepoSecret("TRELLO_BOARD_ID", trelloId.boardId); err != nil {
		return err
	}

	// add todolist id
	if err := c.AddRepoSecret("TRELLO_TODO_LIST_ID", trelloId.todoListId); err != nil {
		return err
	}

	// add doinglist id
	if err := c.AddRepoSecret("TRELLO_DOING_LIST_ID", trelloId.doingListId); err != nil {
		return err
	}

	// add donelist id
	if err := c.AddRepoSecret("TRELLO_DONE_LIST_ID", trelloId.doneListId); err != nil {
		return err
	}

	// add member map
	if err := c.AddRepoSecret("TRELLO_MEMBER_MAP", "[]"); err != nil {
		return err
	}

	return nil
}
