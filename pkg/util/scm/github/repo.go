package github

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/pkgerror"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	webhookName                               = "devstream_webhook"
	defaultLatestCodeZipfileDownloadUrlFormat = "https://codeload.github.com/%s/%s/zip/refs/heads/%s?archive=zip"
)

// DownloadRepo will download repo, return repo local path and error
func (c *Client) DownloadRepo() (string, error) {
	latestCodeZipfileDownloadLocation := downloader.ResourceLocation(fmt.Sprintf(
		defaultLatestCodeZipfileDownloadUrlFormat, c.GetRepoOwner(), c.Repo, c.Branch,
	))
	log.Debugf("github get repo download url: %s.", latestCodeZipfileDownloadLocation)
	log.Info("github start to download repoTemplate...")
	return latestCodeZipfileDownloadLocation.Download()
}

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

func (c *Client) DescribeRepo() (*git.RepoInfo, error) {
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
	repoInfo := &git.RepoInfo{
		Repo:     repo.GetName(),
		Owner:    repo.GetOwner().GetLogin(),
		Org:      repo.GetOrganization().GetLogin(),
		RepoType: "github",
		CloneURL: git.ScmURL(repo.GetCloneURL()),
	}
	return repoInfo, nil
}

func (c *Client) InitRepo() error {
	// It's ok to give the opts.Org to CreateRepo() when create a repository for an authenticated user.
	if err := c.CreateRepo(c.Org, c.Branch); err != nil {
		// recreate if set tryTime
		log.Errorf("Failed to create repo: %s.", err)
		return err
	}
	log.Successf("The repo %s has been created.", c.Repo)

	//	upload a placeholder file to make repo not empty
	if err := c.CreateFile([]byte(" "), repoPlaceHolderFileName, c.Branch); err != nil {
		log.Debugf("Failed to add the first file: %s. ", err)
		return err
	}
	log.Debugf("Added the .placeholder file.")
	return nil
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

	_, _, err = c.Repositories.UpdateBranchProtection(c.Context, repo.Owner, repo.Repo, branch, req)
	if err != nil {
		return err
	}

	log.Infof("The branch \"%s\" has been protected", branch)
	return nil
}

func (c *Client) AddWebhook(webhookConfig *git.WebhookConfig) error {
	hook := new(github.Hook)
	hook.Name = github.String(webhookName)
	hook.Events = []string{"pull_request", "push"}
	hook.Config = map[string]interface{}{}
	hook.Config["url"] = webhookConfig.Address
	hook.Config["content_type"] = "json"
	_, _, err := client.Repositories.CreateHook(c.Context, c.Owner, c.Repo, hook)
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errHookAlreadyExist) {
		return err
	}
	return nil
}

func (c *Client) DeleteWebhook(webhookConfig *git.WebhookConfig) error {
	// list 100 webhooks of this repo
	allHooks, _, err := client.Repositories.ListHooks(
		c.Context, c.Owner, c.Repo, &github.ListOptions{
			PerPage: 100,
		},
	)
	if err != nil {
		log.Debugf("github list webhook failed: %v", err)
		return err
	}
	for _, hook := range allHooks {
		if *hook.Name == webhookName {
			_, err := client.Repositories.DeleteHook(c.Context, c.Owner, c.Repo, *hook.ID)
			return err
		}
	}
	log.Debugf("github webhook is already deleted")
	return nil
}
