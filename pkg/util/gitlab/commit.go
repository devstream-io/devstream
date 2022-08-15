package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// https://docs.gitlab.com/ee/api/commits.html
// https://github.com/xanzy/go-gitlab/blob/master/commits.go
func (c *Client) DeleteFiles(commitInfo *git.CommitInfo) error {
	deleteCommitoptions := c.CreateCommitInfo(gitlab.FileDelete, commitInfo)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), deleteCommitoptions)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateFiles(commitInfo *git.CommitInfo) error {
	updateCommitoptions := c.CreateCommitInfo(gitlab.FileUpdate, commitInfo)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), updateCommitoptions)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PushLocalFileToRepo(commitInfo *git.CommitInfo) (bool, error) {
	createCommitOptions := c.CreateCommitInfo(gitlab.FileCreate, commitInfo)
	_, response, err := c.Commits.CreateCommit(c.GetRepoPath(), createCommitOptions)
	log.Debug(response.Body)
	if err != nil {
		return true, err
	}
	return false, nil
}

func (c *Client) CreateCommitInfo(action gitlab.FileActionValue, commitInfo *git.CommitInfo) *gitlab.CreateCommitOptions {
	var commitActionsOptions = make([]*gitlab.CommitActionOptions, 0, len(commitInfo.GitFileMap))

	for fileName, content := range commitInfo.GitFileMap {
		commitActionsOptions = append(commitActionsOptions, &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(action),
			FilePath: gitlab.String(fileName),
			Content:  gitlab.String(string(content)),
		})
	}
	return &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(c.Branch),
		CommitMessage: gitlab.String(commitInfo.CommitMsg),
		Actions:       commitActionsOptions,
	}
}
