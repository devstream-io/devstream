package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func (c *Client) GetLastCommit() (*github.RepositoryCommit, error) {
	commits, _, err := c.Client.Repositories.ListCommits(c.Context, c.GetRepoOwner(), c.Repo, &github.CommitsListOptions{})
	if err != nil {
		log.Debugf("Failed to get RepositoryCommits: %s.", err)
		return nil, err
	}

	if len(commits) == 0 {
		msg := "no commits was found"
		log.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	return commits[0], nil
}

func (c *Client) BuildCommitTree(ref *github.Reference, commitInfo *git.CommitInfo, checkChange bool) (*github.Tree, error) {
	var entries []*github.TreeEntry
	for githubPath, content := range commitInfo.GitFileMap {
		if checkChange && !c.checkFileChange(githubPath, content) {
			log.Debugf("Github File [%s] content not changed, not commit", githubPath)
			continue
		}
		entries = append(entries, &github.TreeEntry{
			Path:    github.String(githubPath),
			Type:    github.String("blob"),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		})
	}
	if len(entries) == 0 {
		return nil, nil
	}
	tree, _, err := client.Git.CreateTree(c.Context, c.GetRepoOwner(), c.Repo, *ref.Object.SHA, entries)
	return tree, err
}
