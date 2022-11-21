package github

import (
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

// PushFiles will push local change to remote repo
// return boolean value is for control whether to roll out if encounter error
func (c *Client) PushFiles(commitInfo *git.CommitInfo, checkChange bool) (bool, error) {
	// 1. create new branch from main
	ref, err := c.NewBranch(commitInfo.CommitBranch)
	if err != nil {
		log.Warnf("Failed to create transit branch: %s", err)
		return false, err
	}
	// delete new branch after func exit
	defer func() {
		err = c.DeleteBranch(commitInfo.CommitBranch)
		if err != nil {
			log.Warnf("Failed to delete transit branch: %s", err)
		}
	}()
	tree, err := c.BuildCommitTree(ref, commitInfo, checkChange)
	if err != nil {
		log.Debugf("Failed to build commit tree: %s.", err)
		return true, err
	}
	// if no new files to commit, just return
	if tree == nil {
		log.Successf("Github file all not change, pass...")
		return false, nil
	}

	// 2. push local file change to new branch
	if err := c.PushLocalPath(ref, tree, commitInfo); err != nil {
		log.Debugf("Failed to walk local repo-path: %s.", err)
		return true, err
	}

	// 3. merge new branch to main
	if err = c.MergeCommits(commitInfo); err != nil {
		log.Debugf("Failed to merge commits: %s.", err)
		return true, err
	}
	// 4. delete placeholder file
	if err = c.DeleteFiles(&git.CommitInfo{
		CommitMsg:    "delete placeholder file",
		CommitBranch: c.Branch,
		GitFileMap: git.GitFileContentMap{
			repoPlaceHolderFileName: []byte{},
		},
	}); err != nil {
		log.Debugf("github delete init file failed: %s", err)
	}
	return false, nil
}

func (c *Client) CreateFile(content []byte, filePath, targetBranch string) error {
	defaultMsg := "Initialize the repository"

	opt := &github.RepositoryContentFileOptions{
		Message: &defaultMsg,
		Content: content,
		Branch:  &targetBranch,
	}

	_, _, err := c.Repositories.CreateFile(c.Context, c.GetRepoOwner(), c.Repo, filePath, opt)
	return err
}

func (c *Client) PushLocalPath(ref *github.Reference, tree *github.Tree, commitInfo *git.CommitInfo) error {
	parent, _, err := client.Repositories.GetCommit(c.Context, c.GetRepoOwner(), c.Repo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: github.String(defaultCommitAuthor), Email: github.String(defaultCommitAuthorEmail)}
	commit := &github.Commit{Author: author, Message: github.String(commitInfo.CommitMsg), Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(c.Context, c.GetRepoOwner(), c.Repo, commit)
	if err != nil {
		return err
	}
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(c.Context, c.GetRepoOwner(), c.Repo, ref, false)
	return err
}

func (c *Client) GetPathInfo(location string) ([]*git.RepoFileStatus, error) {
	fileContent, directoryContent, resp, err := c.Client.Repositories.GetContents(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		location,
		&github.RepositoryContentGetOptions{Ref: c.Branch},
	)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	gitfilesStatus := make([]*git.RepoFileStatus, 0, len(directoryContent)+1)
	if fileContent != nil {
		gitfilesStatus = append(gitfilesStatus, buildGitFileInfoFromRepContent(c.Branch, fileContent))
	}
	for _, repFileContent := range directoryContent {
		gitfilesStatus = append(gitfilesStatus, buildGitFileInfoFromRepContent(c.Branch, repFileContent))
	}
	return gitfilesStatus, nil
}

func (c *Client) checkFileChange(location string, content []byte) bool {
	fileInfos, err := c.GetPathInfo(location)
	if err != nil {
		log.Debugf("Github request check file SHA failed: %s", err)
		return true
	}
	contentSHA := git.CalculateGitHubBlobSHA(content)
	for _, f := range fileInfos {
		if f.SHA == contentSHA {
			return false
		}
	}
	return true
}

func (c *Client) DeleteFiles(commitInfo *git.CommitInfo) error {
	for fileLoc := range commitInfo.GitFileMap {
		fileInfos, err := c.GetPathInfo(fileLoc)
		if err != nil || len(fileInfos) != 1 {
			log.Debugf("Github file %s already removed.", fileLoc)
			continue
		}
		opts := fileInfos[0].EncodeToGitHubContentOption(
			commitInfo.CommitMsg,
		)
		log.Debugf("Deleting GitHub file %s ...", fileLoc)
		_, _, err = c.Client.Repositories.DeleteFile(
			c.Context,
			c.GetRepoOwner(),
			c.Repo,
			fileLoc,
			opts)
		if err != nil {
			return err
		}
		log.Debugf("GitHub file %s removed.", fileLoc)
	}
	return nil
}

func buildGitFileInfoFromRepContent(branch string, repContent *github.RepositoryContent) *git.RepoFileStatus {
	return &git.RepoFileStatus{
		Path:   *repContent.Path,
		SHA:    *repContent.SHA,
		Branch: branch,
	}
}
