package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) NewBranch(newBranch string) (*github.Reference, error) {
	ref, err := c.getMainBranchRef(c.Branch)
	if err != nil {
		return nil, err
	}

	newRef := fmt.Sprintf("heads/%s", newBranch)
	newRefItem, _, err := c.Git.CreateRef(c.Context, c.GetRepoOwner(), c.Repo, &github.Reference{
		Ref: &newRef,
		Object: &github.GitObject{
			Type: nil,
			SHA:  ref.GetObject().SHA,
			URL:  nil,
		},
	})
	return newRefItem, err
}

func (c *Client) DeleteBranch(branch string) error {
	refStr := fmt.Sprintf("heads/%s", branch)
	log.Debugf("Deleting ref: %s.", refStr)

	_, err := c.Git.DeleteRef(c.Context, c.GetRepoOwner(), c.Repo, refStr)
	return err
}

func (c *Client) MergeCommits(commitInfo *git.CommitInfo) error {
	number, err := c.NewPullRequest(commitInfo)
	if err != nil {
		return err
	}

	return c.MergePullRequest(number, MergeMethodMerge)
}

func (c *Client) getMainBranchRef(branch string) (*github.Reference, error) {
	refStr := fmt.Sprintf("heads/%s", branch)
	ref, _, err := c.Git.GetRef(c.Context, c.GetRepoOwner(), c.Repo, refStr)
	if err != nil {
		log.Debugf("Failed to get the ref for %s: %s.", refStr, err)
		return nil, err
	}
	log.Debugf("Got the ref: Ref %s, JenkinsURL %s, nodeId %s, Obj: %s.",
		ref.GetRef(), ref.GetURL(), ref.GetNodeID(), ref.GetObject().String())
	return ref, nil
}
