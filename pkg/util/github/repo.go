package github

import (
	"net/http"
	"strings"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) CreateRepo(org, defaultBranch string) error {
	repo := &github.Repository{
		Name:          &c.Repo,
		DefaultBranch: github.String(defaultBranch),
	}

	if org != "" {
		log.Infof("Prepare to create an organization repository: %s/%s", org, repo.GetName())
	}
	_, _, err := c.Repositories.Create(c.Context, org, repo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteRepo() error {
	response, err := c.Client.Repositories.Delete(c.Context, c.GetRepoOwner(), c.Repo)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		log.Errorf("Delete repo failed: %s.", err)
		return err
	}

	if response.StatusCode == http.StatusNotFound {
		log.Infof("GitHub repo %s was not found. Nothing to does here.", c.Repo)
		return nil
	}

	log.Successf("GitHub repo %s removed.", c.Repo)
	return nil
}

func (c *Client) DescribeRepo() (*github.Repository, error) {
	repo, resp, err := c.Client.Repositories.Get(
		c.Context,
		c.GetRepoOwner(),
		c.Repo)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// PushLocalPathToBranch will push local change to remote repo
// return boolean value is for control whether to rollout if encounter error
func (c *Client) PushLocalFileToRepo(commitInfo *git.CommitInfo) (bool, error) {
	// 1. create new branch from main
	ref, err := c.NewBranch(commitInfo.CommitBranch)
	if err != nil {
		log.Debugf("Failed to create transit branch: %s", err)
		return false, err
	}
	// delete new branch after func exit
	defer func() {
		err = c.DeleteBranch(commitInfo.CommitBranch)
		if err != nil {
			log.Warnf("Failed to delete transit branch: %s", err)
		}
	}()
	tree, err := c.BuildCommitTree(ref, commitInfo)

	// 2. push local file change to new branch
	if err := c.PushLocalPath(ref, tree, commitInfo); err != nil {
		log.Debugf("Failed to walk local repo-path: %s.", err)
		return true, err
	}

	// 3. merge new branch to main
	if err = c.MergeCommits(commitInfo.CommitBranch); err != nil {
		log.Debugf("Failed to merge commits: %s.", err)
		return true, err
	}
	return false, nil
}

func (c *Client) InitRepo() error {
	// It's ok to give the opts.Org to CreateRepo() when create a repository for a authenticated user.
	if err := c.CreateRepo(c.Org, c.Branch); err != nil {
		// recreate if set tryTime
		log.Errorf("Failed to create repo: %s.", err)
		return err
	}
	log.Infof("The repo %s has been created.", c.Repo)

	// upload a placeholder file to make repo not empty
	if err := c.CreateFile([]byte(" "), ".placeholder", c.Branch); err != nil {
		log.Debugf("Failed to add the first file: %s.", err)
		return err
	}
	log.Debugf("Added the .placeholder file.")
	return nil
}

func (c *Client) PushInitRepo(commitInfo *git.CommitInfo) error {
	// 1. init repo
	if err := c.InitRepo(); err != nil {
		return err
	}

	// if encounter rollout error, delete repo
	var needRollBack bool
	defer func() {
		if !needRollBack {
			return
		}
		// need to clean the repo created when retErr != nil
		if err := c.DeleteRepo(); err != nil {
			log.Errorf("Failed to delete the repo %s: %s.", c.Repo, err)
		}
	}()

	// 2. push local path to repo
	needRollBack, err := c.PushLocalFileToRepo(commitInfo)
	if err != nil {
		return err
	}

	// 3. protect branch
	err = c.ProtectBranch(c.Branch)
	return err
}

// ProtectBranch will protect the special branch
func (c *Client) ProtectBranch(branch string) error {
	req := &github.ProtectionRequest{
		EnforceAdmins: false,
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			RequireCodeOwnerReviews:      true,
			DismissStaleReviews:          true,
			RequiredApprovingReviewCount: 1,
		},
		RequiredConversationResolution: github.Bool(true),
	}

	repo, err := c.DescribeRepo()
	if err != nil {
		return err
	}

	_, _, err = c.Repositories.UpdateBranchProtection(c.Context, repo.GetOwner().GetLogin(), repo.GetName(), branch, req)
	if err != nil {
		return err
	}

	log.Infof("The branch \"%s\" has been protected", branch)
	return nil
}
