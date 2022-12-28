package scm

import (
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

func NewClientWithAuth(repoInfo *git.RepoInfo) (ClientOperation, error) {
	repoInfo.NeedAuth = true
	return NewClient(repoInfo)
}

func NewClient(repoInfo *git.RepoInfo) (ClientOperation, error) {
	if err := repoInfo.SetDefault(); err != nil {
		return nil, err
	}
	switch repoInfo.RepoType {
	case "github":
		return github.NewClient(repoInfo)
	default:
		// default use gitlab repo
		return gitlab.NewClient(repoInfo)
	}
}

type ClientOperation interface {
	InitRepo() error
	DeleteRepo() error
	DownloadRepo() (string, error)
	DescribeRepo() (*git.RepoInfo, error)
	PushFiles(commitInfo *git.CommitInfo, checkUpdate bool) (bool, error)
	DeleteFiles(commitInfo *git.CommitInfo) error
	GetPathInfo(path string) ([]*git.RepoFileStatus, error)
	AddWebhook(webhookConfig *git.WebhookConfig) error
	DeleteWebhook(webhookConfig *git.WebhookConfig) error
	AddRepoSecret(secretKey, secretValue string) error
}

func PushInitRepo(client ClientOperation, commitInfo *git.CommitInfo) error {
	// 1. init repo
	if err := client.InitRepo(); err != nil {
		return err
	}

	var (
		// if encounter rollout error, delete repo
		needRollBack bool
		err          error
	)
	defer func() {
		if !needRollBack {
			return
		}
		// need to clean the repo created when reterr != nil
		if err := client.DeleteRepo(); err != nil {
			log.Errorf("failed to delete the repo: %s.", err)
		}
	}()

	// 2. push local path to repo
	needRollBack, err = client.PushFiles(commitInfo, false)
	return err
}
