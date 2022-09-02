package github

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/scm/git"

	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/google/go-github/v42/github"
)

type MergeMethod string

const (
	MergeMethodSquash MergeMethod = "squash"
	MergeMethodMerge  MergeMethod = "merge"
	MergeMethodRebase MergeMethod = "rebase"
)

func (c *Client) NewPullRequest(commitInfo *git.CommitInfo) (int, error) {
	title := commitInfo.CommitMsg
	head := commitInfo.CommitBranch
	base := c.Branch
	body := title
	mcm := false
	draft := false

	pr, _, err := c.PullRequests.Create(c.Context, c.GetRepoOwner(), c.Repo, &github.NewPullRequest{
		Title:               &title,
		Head:                &head,
		Base:                &base,
		Body:                &body,
		Issue:               nil,
		MaintainerCanModify: &mcm,
		Draft:               &draft,
	})

	if err != nil {
		log.Debugf("Failed to create the pr: %s.", err)
		return 0, err
	}

	log.Debugf("The pr has been created: #%d.", pr.GetNumber())

	return pr.GetNumber(), nil
}

func (c *Client) MergePullRequest(number int, mergeMethod MergeMethod) error {
	commitMsg := "Initialized by DevStream"
	ret, _, err := c.PullRequests.Merge(c.Context, c.GetRepoOwner(), c.Repo, number, commitMsg, &github.PullRequestOptions{
		CommitTitle:        commitMsg,
		SHA:                "",
		MergeMethod:        string(mergeMethod),
		DontDefaultIfBlank: false,
	})
	if err != nil {
		log.Debugf("Got an error when merge the pr: %s.", err)
		return err
	}

	if !ret.GetMerged() {
		return fmt.Errorf("merge failed")
	}

	return nil
}
