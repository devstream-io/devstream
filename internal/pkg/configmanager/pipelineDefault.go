package configmanager

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/log"
)

var optionConfiguratorMap = map[string]pipelineOption{
	"github-actions":   githubGeneral,
	"gitlab-ci":        gitlabGeneral,
	"jenkins-pipeline": jenkinsGeneral,
	"argocdapp":        argocdApp,
}

type pipelineOptionGenerator func(originOption RawOptions, app *app) RawOptions

type pipelineOption struct {
	defaultConfigLocation string
	optionGeneratorFunc   pipelineOptionGenerator
}

// TODO(steinliber) unify all ci/cd config to same config options
var (
	// github actions pipeline options
	githubGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/ci-template/main/github-actions/workflows/main.yml",
		optionGeneratorFunc:   pipelineGeneralGenerator,
	}
	gitlabGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/ci-template/main/gitlab-ci/.gitlab-ci.yml",
		optionGeneratorFunc:   pipelineGeneralGenerator,
	}
	jenkinsGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/ci-template/main/jenkins-pipeline/general/Jenkinsfile",
		optionGeneratorFunc:   jenkinsGenerator,
	}
	argocdApp = pipelineOption{
		optionGeneratorFunc: pipelineArgocdAppGenerator,
	}
)

// jenkinsGenerator generate jenkins pipeline config
func jenkinsGenerator(options RawOptions, app *app) RawOptions {
	newOptions := pipelineGeneralGenerator(options, app)
	// extract jenkins config from options
	jenkinsOptions, exist := options["jenkins"]
	if exist {
		newOptions["jenkins"] = jenkinsOptions
	}
	return newOptions
}

// pipelineGeneralGenerator generate pipeline general options from RawOptions
func pipelineGeneralGenerator(options RawOptions, app *app) RawOptions {
	if app.Spec != nil {
		app.Spec.updatePiplineOption(options)
	}
	// update image related config
	newOption := make(RawOptions)
	newOption["pipeline"] = options
	newOption["scm"] = RawOptions(app.Repo.Encode())
	return newOption
}

// pipelineArgocdAppGenerator generate argocdApp options from RawOptions
func pipelineArgocdAppGenerator(options RawOptions, app *app) RawOptions {
	// config app default options
	if _, exist := options["app"]; !exist {
		options["app"] = RawOptions{
			"name":      app.Name,
			"namespace": "argocd",
		}
	}
	// config destination options
	if _, exist := options["destination"]; !exist {
		options["destination"] = RawOptions{
			"server":    "https://kubernetes.default.svc",
			"namespace": "default",
		}
	}
	// config source default options
	repoInfo, err := app.Repo.BuildRepoInfo()
	if err != nil {
		log.Errorf("parse argocd repoInfo failed: %+v", err)
		return options
	}
	if source, sourceExist := options["source"]; sourceExist {
		sourceMap := source.(RawOptions)
		if _, repoURLExist := sourceMap["repoURL"]; !repoURLExist {
			sourceMap["repoURL"] = repoInfo.CloneURL
		}
		options["source"] = sourceMap
	} else {
		options["source"] = RawOptions{
			"valuefile": "values.yaml",
			"path":      fmt.Sprintf("helm/%s", app.Name),
			"repoURL":   repoInfo.CloneURL,
		}
	}
	return options
}

// hasDefaultConfig check whether
func (o *pipelineOption) hasDefaultConfig() bool {
	return o.defaultConfigLocation != ""
}
