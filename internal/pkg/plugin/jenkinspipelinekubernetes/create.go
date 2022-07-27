package jenkinspipelinekubernetes

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validateAndHandleOptions(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(&opts)
	if err != nil {
		return nil, err
	}

	// create credential if not exists
	if _, err := client.GetCredentialsUsername(jenkinsCredentialID); err != nil {
		log.Infof("credential %s not found, creating...", jenkinsCredentialID)
		if err := client.CreateCredentialsUsername(jenkinsCredentialUsername, opts.GitHubToken, jenkinsCredentialID, jenkinsCredentialDesc); err != nil {
			return nil, err
		}
	}

	// create job if not exists
	ctx := context.Background()
	if _, err := client.GetJob(ctx, opts.J.JobName); err != nil {
		log.Infof("job %s not found, creating...", opts.J.JobName)
		jobXmlOpts := &JobXmlOptions{
			GitHubRepoURL:      opts.GitHubRepoURL,
			CredentialsID:      jenkinsCredentialID,
			PipelineScriptPath: opts.J.PipelineScriptPath,
		}
		jobXmlContent := renderJobXml(jobTemplate, jobXmlOpts)
		if _, err := client.CreateJob(context.Background(), jobXmlContent, opts.J.JobName); err != nil {
			return nil, fmt.Errorf("failed to create job: %s", err)
		}
	}

	// TODO(aFlyBird0): use JCasC to configure job creation
	// configuration as code: update jenkins config
	//        + configure system -> GitHub Pull Request Builder -> Jenkins URL override
	//        + configure system -> GitHub Pull Request Builder -> git personal access token
	// - job creation:
	//        + pr builder
	//        + main builder
	//        + GitHub project Project url
	//        + select GitHub Pull Request Builder under Build Triggers
	//        + Trigger phrase: optional
	//        + Pipeline Definition: pipeline script from SCM
	//        + Branch Specifier & Refspec: PR/main
	//        + unselect Lightweight checkout
	//        + Jenkinsfile template -> github repo (https://github.com/IronCore864/jenkins-test/blob/main/Jenkinsfile, https://plugins.jenkins.io/kubernetes/)

	// TODO(aFlyBird0): about how to create an new config yaml in JCasC:
	// refer: https://plugins.jenkins.io/configuration-as-code/
	// refer: https://github.com/jenkinsci/helm-charts/blob/e4242561e5ea205bfa3405064cf5fe39b5a22d93/charts/jenkins/values.yaml#L341
	// refer: https://github.com/jenkinsci/helm-charts/blob/e4242561e5ea205bfa3405064cf5fe39b5a22d93/charts/jenkins/templates/jcasc-config.yaml#L1-L25
	// example: If we want to create a new config yaml in JCasC, we can use the following code:
	// 			1. render this ConfigMap by yourself: https://github.com/jenkinsci/helm-charts/blob/e4242561e5ea205bfa3405064cf5fe39b5a22d93/charts/jenkins/templates/jcasc-config.yaml#L6-L26
	//          2. $key should be your config file name(without .yaml extension), value should be the real config file content
	//			3. once you want to put a new config yaml into $JenkinsHome/config-as-code/ folder, just use k8s sdk to create a new ConfigMap
	//			4. don't forget to add a label to the ConfigMap, such as "author: devstream". So that we could easy to filter the ConfigMap created by devstream to delete them.
	// 			5. if you want to update an existing config yaml, just use k8s sdk to update the ConfigMap
	//			6. it seems that once you create/update a ConfigMap, Jenkins(installed by helm) will automatically reload the config yaml,
	//				if not, you can use the following api to reload the config yaml: "POST JENKINS_URL/configuration-as-code/reload"（ remember to add jenkins crumb in http header)
	//				see here for info:https://github.com/jenkinsci/configuration-as-code-plugin/blob/master/docs/features/configurationReload.md
	// there are many things to do:
	// 1. figure out the JCasC content we want to create
	// 2. create a new ConfigMap according to the explanation above
	// 3. handle update of ConfigMap
	// 4. add "read" functions to the ConfigMap to get the content of the ConfigMap and check if the resource is drifted
	// 5. maybe we also should consider to expose some config key in ConfigMap to the user

	// TODO(aFlyBird0): what if user doesn't use helm to install jenkins? then the JCasC may not be automatically reloaded.

	res := &resource{
		CredentialsCreated: true,
		JobCreated:         true,
	}

	return res.toMap(), nil
}

type JobXmlOptions struct {
	GitHubRepoURL      string
	CredentialsID      string
	PipelineScriptPath string
}

func renderJobXml(jobTemplate string, opts *JobXmlOptions) string {
	// TODO(aFlyBird0): use html/template to generate the job template. It's a good first issue. :)
	replacerSlice := []string{
		"{{.GitHubRepoURL}}", opts.GitHubRepoURL,
		"{{.CredentialsID}}", opts.CredentialsID,
		"{{.PipelineScriptPath}}", opts.PipelineScriptPath,
	}

	jobXml := strings.NewReplacer(replacerSlice...).Replace(jobTemplate)

	log.Debugf("job xml rendered: %s", jobXml)

	return jobXml
}
