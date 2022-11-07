package gitlab

import (
	"net/http"

	"github.com/xanzy/go-gitlab"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func (c *Client) PushFiles(commitInfo *git.CommitInfo, checkUpdate bool) (bool, error) {
	// if checkUpdate is true, check files to update first
	tree := newCommitTree(commitInfo.CommitMsg, c.Branch)
	if checkUpdate {
		for scmPath, content := range commitInfo.GitFileMap {
			fileExist, err := c.checkFileExist(scmPath)
			if err != nil {
				return false, err
			}
			if fileExist {
				tree.addCommitFile(gitlab.FileUpdate, scmPath, content)
			} else {
				tree.addCommitFile(gitlab.FileCreate, scmPath, content)
			}
		}
	} else {
		tree.addCommitFilesFromMap(gitlab.FileCreate, commitInfo.GitFileMap)
	}
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), tree.createCommitInfo())
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errFileExist) {
		return true, c.newModuleError(err)
	}
	return false, nil
}

func (c *Client) DeleteFiles(commitInfo *git.CommitInfo) error {
	tree := newCommitTree(commitInfo.CommitMsg, c.Branch)
	tree.addCommitFilesFromMap(gitlab.FileDelete, commitInfo.GitFileMap)
	_, _, err := c.Commits.CreateCommit(c.GetRepoPath(), tree.createCommitInfo())
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errRepoNotFound) {
		return c.newModuleError(err)
	}
	return nil
}

func (c *Client) checkFileExist(filename string) (bool, error) {
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

func (c *Client) GetPathInfo(path string) ([]*git.RepoFileStatus, error) {
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
