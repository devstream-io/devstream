package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func (c *Client) PushLocalFileToRepo(commitInfo *git.CommitInfo, checkUpdate bool) (bool, error) {
	// if checkUpdate is true, check files to update first
	if checkUpdate {
		updateCommitInfo := &git.CommitInfo{
			CommitMsg:    commitInfo.CommitMsg,
			CommitBranch: commitInfo.CommitBranch,
		}
		commitMap := make(git.GitFileContentMap)
		for filePath, content := range commitInfo.GitFileMap {
			fileExist, err := c.FileExists(filePath)
			if err != nil {
				return false, err
			}
			if fileExist {
				commitMap[filePath] = content
				delete(commitInfo.GitFileMap, filePath)
			}
		}
		if len(commitMap) > 0 {
			err := c.UpdateFiles(updateCommitInfo)
			if err != nil {
				return true, err
			}
		}
	}
	createCommitOptions := c.CreateCommitInfo(gitlab.FileCreate, commitInfo)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), createCommitOptions)
	if err != nil && pkgerror.CheckSlientErrorByMessage(err, errFileExist) {
		return true, c.newModuleError(err)
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
