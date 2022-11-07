package jenkinspipeline

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/random"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

func installPipeline(options configmanager.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. init jenkins Client
	jenkinsClient, err := opts.Jenkins.newClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	// 2. generate secretToken for webhook auth
	secretToken := random.GenerateRandomSecretToken()
	if err := opts.install(jenkinsClient, secretToken); err != nil {
		log.Debugf("jenkins install pipeline failed: %s", err)
		return err
	}
	// 3. create or update scm webhook
	scmClient, err := scm.NewClientWithAuth(opts.ProjectRepo)
	if err != nil {
		return err
	}
	webHookConfig := opts.ProjectRepo.BuildWebhookInfo(
		opts.Jenkins.URL, string(opts.JobName), secretToken,
	)
	return scmClient.AddWebhook(webHookConfig)

}

func deletePipeline(options configmanager.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. delete jenkins job
	client, err := opts.Jenkins.newClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	err = opts.remove(client)
	if err != nil {
		return err
	}
	// 2. delete repo webhook
	scmClient, err := scm.NewClientWithAuth(opts.ProjectRepo)
	if err != nil {
		return err
	}
	webHookConfig := opts.ProjectRepo.BuildWebhookInfo(
		opts.Jenkins.URL, string(opts.JobName), "",
	)
	return scmClient.DeleteWebhook(webHookConfig)
}
