package jenkinspipelinekubernetes

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	jenkinsCredentialID   = "credential-jenkins-pipeline-kubernetes-by-devstream"
	jenkinsCredentialDesc = "Jenkins Pipeline secret, created by devstream/jenkins-pipeline-kubernetes"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	client := NewClient(&opts)

	// always try to create credential, there will be no error if it already exists
	if err := client.CreateCredentialUsernamePassword(); err != nil {
		return nil, err
	}

	jobXmlContent := renderJobXml(jobTemplate, &opts)

	// TODO(aFlyBird0): check if the job already exists

	if err := client.CreateItem(jobXmlContent); err != nil {
		return nil, fmt.Errorf("failed to create job: %s", err)
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

	// TODO(aFlyBird0): build dtm resource

	// TODO(aFlyBird0): what if user doesn't use helm to install jenkins? then the JCasC may not be automatically reloaded.

	return (&resource{}).toMap(), nil
}

// TODO(aFlyBird0): unit test
// TODO(aFlyBird0): now jenkins script path is hardcoded to "Jenkinsfile", it should be configurable
func renderJobXml(jobTemplate string, opts *Options) string {
	// note: maybe it is better to use html/template to generate the job template,
	// but that way is complex and this is the simplest way to do it
	jobXml := strings.Replace(jobTemplate, "{{.GitHubRepoURL}}", opts.GitHubRepoURL, 1)
	jobXml = strings.Replace(jobXml, "{{.CredentialsID}}", jenkinsCredentialID, 1)

	log.Debugf("job xml rendered: %s", jobXml)
	return jobXml
}
