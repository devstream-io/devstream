package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

func CreateOrUpdateJob(options plugininstaller.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. init repo webhook
	jenkinsClient, err := opts.newJenkinsClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	// 2. create or update jenkins job
	return opts.createOrUpdateJob(jenkinsClient)
}

func CreateRepoWebhook(options plugininstaller.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	scmClient, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	return scmClient.AddWebhook(opts.buildWebhookInfo())
}

func DeleteJob(options plugininstaller.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. delete jenkins job
	client, err := opts.newJenkinsClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	err = opts.deleteJob(client)
	if err != nil {
		return err
	}
	// 2. delete repo webhook
	scmClient, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	return scmClient.DeleteWebhook(opts.buildWebhookInfo())
}

// PreInstall will download jenkins plugins and config jenkins casc
func PreInstall(options plugininstaller.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. init jenkins client
	jenkinsClient, err := opts.newJenkinsClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}

	// 2. get all plugins need to preConfig
	pluginsConfigs := opts.extractJenkinsPlugins()

	// 3. install all plugins
	err = installPlugins(jenkinsClient, pluginsConfigs, opts.Jenkins.EnableRestart)
	if err != nil {
		log.Debugf("jenkins preinstall plugins failed: %s", err)
		return err
	}

	// 4. config repo related config in jenkins
	return preConfigPlugins(jenkinsClient, pluginsConfigs)
}
