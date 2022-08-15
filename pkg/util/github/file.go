package github

import (
	"time"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/git"
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
