package git

import (
	"github.com/devstream-io/devstream/pkg/util/log"
)

type ClientOperation interface {
	InitRepo() error
	DeleteRepo() error
	PushLocalFileToRepo(commitInfo *CommitInfo, checkUpdate bool) (bool, error)
	GetLocationInfo(path string) ([]*RepoFileStatus, error)
	DeleteFiles(commitInfo *CommitInfo) error
	AddWebhook(webhookConfig *WebhookConfig) error
	DeleteWebhook(webhookConfig *WebhookConfig) error
}

func PushInitRepo(client ClientOperation, commitInfo *CommitInfo) error {
	// 1. init repo
	if err := client.InitRepo(); err != nil {
		return err
	}

	// if encounter rollout error, delete repo
	var needRollBack bool
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
	needRollBack, err := client.PushLocalFileToRepo(commitInfo, false)
	return err
}
