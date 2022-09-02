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
	err = opts.createOrUpdateJob(jenkinsClient)
	if err != nil {
		log.Debugf("jenkins execute script failed: %s", err)
		return err
	}
	// 3. create repo webhook
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
	client, err := opts.newJenkinsClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	err = opts.deleteJob(client)
	if err != nil {
		return err
	}
	// delete repo webhook
	scmClient, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	return scmClient.DeleteWebhook(opts.buildWebhookInfo())
}

func PreInstall(plugins []string, cascTemplate string) plugininstaller.BaseOperation {
	return func(options plugininstaller.RawOptions) error {
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
		// 2. install plugins
		err = opts.installPlugins(jenkinsClient, plugins)
		if err != nil {
			log.Debugf("jenkins preinstall plugins failed: %s", err)
			return err
		}

		switch opts.ProjectRepo.RepoType {
		case "gitlab":
			// 3. create gitlab connection for gitlab
			return opts.createGitlabConnection(jenkinsClient, cascTemplate)
		default:
			log.Debugf("jenkins preinstall only support gitlab for now")
			return nil
		}
	}
}
