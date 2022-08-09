package github

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/go-github/v42/github"

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
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	response, err := c.Client.Repositories.Delete(c.Context, owner, c.Repo)

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

func (c *Client) GetRepoDescription() (*github.Repository, error) {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	repo, resp, err := c.Client.Repositories.Get(
		c.Context,
		owner,
		c.Repo)

	if err != nil {
		return nil, err
	}

	if repo == nil && resp.StatusCode == http.StatusNotFound {
		return repo, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return repo, nil
}

// PushLocalPathToBranch will push local change to remote repo
// return boolean value is for control whether to rollout if encounter error
func (c *Client) PushLocalPathToBranch(mergeBranch, mainBranch, repoPath string) (bool, error) {
	// 1. create new branch from main
	err := c.NewBranch(mainBranch, mergeBranch)
	if err != nil {
		log.Debugf("Failed to create transit branch: %s", err)
		return false, err
	}
	// 2. push local file change to new branch
	if err := c.PushLocalPath(repoPath, mergeBranch); err != nil {
		log.Debugf("Failed to walk local repo-path: %s.", err)
		return true, err
	}
	// 3. merge new branch to main
	if err = c.MergeCommits(mergeBranch, mainBranch); err != nil {
		log.Debugf("Failed to merge commits: %s.", err)
		return true, err
	}
	// 4. delete new branch
	err = c.DeleteBranch(mergeBranch)
	if err != nil {
		log.Debugf("Failed to delete transit branch: %s", err)
		return false, err
	}
	return false, nil
}

func (c *Client) InitRepo(mainBranch string) error {
	// It's ok to give the opts.Org to CreateRepo() when create a repository for a authenticated user.
	if err := c.CreateRepo(c.Org, mainBranch); err != nil {
		// recreate if set tryTime
		log.Errorf("Failed to create repo: %s.", err)
		return err
	}
	log.Infof("The repo %s has been created.", c.Repo)

	// upload a placeholder file to make repo not empty
	if err := c.CreateFile([]byte(" "), ".placeholder", mainBranch); err != nil {
		log.Debugf("Failed to add the first file: %s.", err)
		return err
	}
	log.Debugf("Added the .placeholder file.")
	return nil
}

func (c *Client) PushInitRepo(transitBranch, branch, localPath string) error {
	// 1. init repo
	if err := c.InitRepo(branch); err != nil {
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
	needRollBack, err := c.PushLocalPathToBranch(transitBranch, branch, localPath)
	if err != nil {
		return err
	}

	// 3. protect branch
	err = c.ProtectBranch(branch)
	return err
}

// ProtectBranch will protect the special branch
func (c *Client) ProtectBranch(branch string) error {
	req := &github.ProtectionRequest{
		RequiredStatusChecks: nil,
		EnforceAdmins:        false,
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			RequireCodeOwnerReviews: true,
			DismissStaleReviews:     true,
		},
		Restrictions: nil,
	}

	repo, err := c.GetRepoDescription()
	if err != nil {
		log.Errorf("Get repo failed: %s.", err)
		return err
	}

	_, response, err := c.Repositories.UpdateBranchProtection(c.Context, repo.GetOwner().GetLogin(), repo.GetName(), branch, req)
	if err != nil {
		log.Errorf("Protect branch failed: %s.", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("Protect branch failed,status code: %d, response: %s.", response.StatusCode, response.Body)
		return errors.New("protect branch failed")
	}

	log.Infof("The branch \"%s\" has been protected", branch)
	return nil
}
