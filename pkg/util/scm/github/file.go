package github

import (
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

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

func (c *Client) GetLocationInfo(location string) ([]*git.RepoFileStatus, error) {
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

func (c *Client) checkFileChange(location string, content string) bool {
	fileInfos, err := c.GetLocationInfo(location)
	if err != nil {
		log.Debugf("Github request check file SHA failed: %s", err)
		return true
	}
	contentSHA := git.CaluateGitHubBlobSHA(content)
	for _, f := range fileInfos {
		if f.SHA == contentSHA {
			return false
		}
	}
	return true
}

func (c *Client) DeleteFiles(commitInfo *git.CommitInfo) error {
	for fileLoc := range commitInfo.GitFileMap {
		fileInfos, err := c.GetLocationInfo(fileLoc)
		if err != nil || len(fileInfos) != 1 {
			log.Successf("Github file %s already removed.", fileLoc)
			continue
		}
		opts := fileInfos[0].EncodeToGitHubContentOption(
			commitInfo.CommitMsg,
		)
		log.Infof("Deleting GitHub  file %s ...", fileLoc)
		_, _, err = c.Client.Repositories.DeleteFile(
			c.Context,
			c.GetRepoOwner(),
			c.Repo,
			fileLoc,
			opts)
		if err != nil {
			return err
		}
		log.Successf("GitHub file %s removed.", fileLoc)
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
