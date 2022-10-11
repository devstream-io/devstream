package gitlab

import (
	"net/http"

	"github.com/xanzy/go-gitlab"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func (c *Client) DeleteFiles(commitInfo *git.CommitInfo) error {
	deleteCommitOptions := c.CreateCommitInfo(gitlab.FileDelete, commitInfo)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), deleteCommitOptions)
	if err != nil && !pkgerror.CheckSlientErrorByMessage(err, errRepoNotFound) {
		return c.newModuleError(err)
	}
	return nil
}

func (c *Client) UpdateFiles(commitInfo *git.CommitInfo) error {
	updateCommitOptions := c.CreateCommitInfo(gitlab.FileUpdate, commitInfo)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), updateCommitOptions)
	if err != nil {
		return c.newModuleError(err)
	}
	return nil
}

func (c *Client) FileExists(filename string) (bool, error) {
	getFileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(c.Branch),
	}

	_, response, err := c.RepositoryFiles.GetFile(c.GetRepoPath(), filename, getFileOptions)
	for _, v := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound} {
		if response.StatusCode == v {
			return false, nil
		}
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) GetLocationInfo(path string) ([]*git.RepoFileStatus, error) {
	gitRepoFileStatus := make([]*git.RepoFileStatus, 0)
	getFileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(c.Branch),
	}

	file, response, err := c.RepositoryFiles.GetFile(c.GetRepoPath(), path, getFileOptions)
	errCodeSet := mapset.NewSet(http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound)

	if response != nil && errCodeSet.Contains(response.StatusCode) {
		return gitRepoFileStatus, nil
	}

	if err != nil {
		return gitRepoFileStatus, err
	}
	gitRepoFileStatus = append(gitRepoFileStatus, &git.RepoFileStatus{
		Path:   file.FilePath,
		Branch: file.Ref,
		SHA:    file.SHA256,
	})
	return gitRepoFileStatus, nil
}
