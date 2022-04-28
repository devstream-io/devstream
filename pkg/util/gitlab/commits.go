package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// https://docs.gitlab.com/ee/api/commits.html
// https://github.com/xanzy/go-gitlab/blob/master/commits.go

func (c *Client) CommitSingleFile(project, branch, commitMessage, filename, content string) error {
	commitActionOptions := gitlab.CommitActionOptions{
		Action:   gitlab.FileAction(gitlab.FileCreate),
		FilePath: gitlab.String(filename),
		Content:  gitlab.String(content),
	}

	createCommitoptions := gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branch),
		CommitMessage: gitlab.String(commitMessage),
		Actions:       []*gitlab.CommitActionOptions{&commitActionOptions},
	}

	_, _, err := c.Commits.CreateCommit(project, &createCommitoptions)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSingleFile(project, branch, commitMessage, filename string) error {
	commitActionOptions := gitlab.CommitActionOptions{
		Action:   gitlab.FileAction(gitlab.FileDelete),
		FilePath: gitlab.String(filename),
	}

	createCommitoptions := gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branch),
		CommitMessage: gitlab.String(commitMessage),
		Actions:       []*gitlab.CommitActionOptions{&commitActionOptions},
	}

	_, _, err := c.Commits.CreateCommit(project, &createCommitoptions)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateSingleFile(project, branch, commitMessage, filename, content string) error {
	commitActionOptions := gitlab.CommitActionOptions{
		Action:   gitlab.FileAction(gitlab.FileUpdate),
		FilePath: gitlab.String(filename),
		Content:  &content,
	}

	createCommitoptions := gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branch),
		CommitMessage: gitlab.String(commitMessage),
		Actions:       []*gitlab.CommitActionOptions{&commitActionOptions},
	}

	_, _, err := c.Commits.CreateCommit(project, &createCommitoptions)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CommitMultipleFiles(project, branch, commitMessage string, files map[string][]byte) error {
	var commitActionsOptions = make([]*gitlab.CommitActionOptions, 0)

	for fileName, content := range files {
		commitActionsOptions = append(commitActionsOptions, &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileCreate),
			FilePath: gitlab.String(fileName),
			Content:  gitlab.String(string(content)),
		})
	}

	createCommitOptions := gitlab.CreateCommitOptions{
		Branch:        gitlab.String(branch),
		CommitMessage: gitlab.String(commitMessage),
		Actions:       commitActionsOptions,
	}

	_, response, err := c.Commits.CreateCommit(project, &createCommitOptions)
	log.Debug(response.Body)
	if err != nil {
		return err
	}

	return nil
}
