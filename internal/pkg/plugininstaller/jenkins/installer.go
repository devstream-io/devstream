package jenkins

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	jenkinsGitlabCredentialName = "gitlabJeninsCredential"
)

func CreateOrUpdateJob(options plugininstaller.RawOptions) error {
	opts, err := newJobOptions(options)
	if err != nil {
		return err
	}
	// 1. render groovy script
	// secretToken is used for webhook secretToken
	secretToken := generateRandomSecretToken()
	renderOptions := opts.buildRenderOptions(secretToken)
	groovyScript, err := jenkins.BuildRenderedScript(renderOptions)
	if err != nil {
		log.Debugf("jenkins redner template failed: %s", err)
		return err
	}
	// 2. execute script to create job
	jenkinsClient, err := opts.newJenkinsClient()
	if err != nil {
		log.Debugf("jenkins init client failed: %s", err)
		return err
	}
	_, err = jenkinsClient.ExecuteScript(groovyScript)
	if err != nil {
		log.Debugf("jenkins execute script failed: %s", err)
		return err
	}
	// 3. create repo webhook
	webhookInfo := opts.buildWebhookInfo(secretToken)
	return opts.ProjectRepo.AddWebHook(webhookInfo)
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
	if _, err = client.GetJob(context.Background(), opts.JobName); err == nil {
		if _, err := client.DeleteJob(context.Background(), opts.JobName); err != nil {
			return err
		}
	}

	// delete repo webhook
	webhookInfo := opts.buildWebhookInfo("")
	return opts.ProjectRepo.DeleteWebhook(webhookInfo)
}

func PreInstall(plugins []string) plugininstaller.BaseOperation {
	return func(options plugininstaller.RawOptions) error {
		opts, err := newJobOptions(options)
		if err != nil {
			return err
		}
		jenkinsClient, err := opts.newJenkinsClient()
		if err != nil {
			log.Debugf("jenkins init client failed: %s", err)
			return err
		}
		// 1. install plugins
		//TODO(steinliber) check plugin install error
		err = jenkinsClient.InstallPluginsIfNotExists(plugins, opts.JenkinsEnableRestart)
		if err != nil {
			log.Debugf("jenkins preinstall plugins failed: %s", err)
			return err
		}
		// 2. config credentials
		err = jenkinsClient.CreateGiltabCredential(jenkinsGitlabCredentialName, os.Getenv("GITLAB_TOKEN"))
		if err != nil {
			log.Debugf("jenkins preinstall credentials failed: %s", err)
			return err
		}
		//TODO(steinliber) use casc to config gitlab connection
		return nil
	}
}

func generateRandomSecretToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:32]
}
