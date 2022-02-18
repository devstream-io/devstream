package gitlab

import (
	"github.com/xanzy/go-gitlab"
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
