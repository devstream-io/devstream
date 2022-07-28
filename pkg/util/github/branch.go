package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) NewBranch(baseBranch, newBranch string) error {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	refStr := fmt.Sprintf("heads/%s", baseBranch)
	ref, _, err := c.Git.GetRef(c.Context, owner, c.Repo, refStr)
	if err != nil {
		log.Debugf("Failed to get the ref for %s: %s.", refStr, err)
		return err
	}
	log.Debugf("Got the ref: Ref %s, URL %s, nodeId %s, Obj: %s.",
		ref.GetRef(), ref.GetURL(), ref.GetNodeID(), ref.GetObject().String())

	newRef := fmt.Sprintf("heads/%s", newBranch)
	_, _, err = c.Git.CreateRef(c.Context, owner, c.Repo, &github.Reference{
		Ref: &newRef,
		Object: &github.GitObject{
			Type: nil,
			SHA:  ref.GetObject().SHA,
			URL:  nil,
		},
	})
	return err
}

func (c *Client) DeleteBranch(branch string) error {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	refStr := fmt.Sprintf("heads/%s", branch)
	log.Debugf("Deleting ref: %s.", refStr)

	_, err := c.Git.DeleteRef(c.Context, owner, c.Repo, refStr)
	return err
}

func (c *Client) MergeCommits(mergeBranch, mainBranch string) error {
	number, err := c.NewPullRequest(mergeBranch, mainBranch)
	if err != nil {
		return err
	}

	return c.MergePullRequest(number, MergeMethodMerge)
}
