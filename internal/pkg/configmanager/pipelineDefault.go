package configmanager

import (
	"fmt"
)

var optionConfiguratorMap = map[string]pipelineOption{
	"github-actions":   githubGeneral,
	"gitlab-ci":        gitlabGeneral,
	"jenkins-pipeline": jenkinsGeneral,
	"argocdapp":        argocdApp,
}

type pipelineOptionGenerator func(originOption RawOptions, vars *pipelineGlobalOption) RawOptions

// pipelineOption is used for config different pipeline type's default config
type pipelineOption struct {
	defaultConfigLocation string
	optionGeneratorFunc   pipelineOptionGenerator
}

// TODO(steinliber) unify all ci/cd config to same config options
var (
	// github actions pipeline options
	githubGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/github-actions/workflows/main.yml",
		optionGeneratorFunc:   pipelineGeneralGenerator,
	}
	gitlabGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/gitlab-ci/.gitlab-ci.yml",
		optionGeneratorFunc:   pipelineGeneralGenerator,
	}
	jenkinsGeneral = pipelineOption{
		defaultConfigLocation: "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/jenkins-pipeline/general/Jenkinsfile",
		optionGeneratorFunc:   jenkinsGenerator,
	}
	argocdApp = pipelineOption{
		optionGeneratorFunc: pipelineArgocdAppGenerator,
	}
)

// jenkinsGenerator generate jenkins pipeline config
func jenkinsGenerator(options RawOptions, globalVars *pipelineGlobalOption) RawOptions {
	newOptions := pipelineGeneralGenerator(options, globalVars)
	// extract jenkins config from options
	jenkinsOptions, exist := options["jenkins"]
	if exist {
		newOptions["jenkins"] = jenkinsOptions
	}
	return newOptions
}

// pipelineGeneralGenerator generate pipeline general options from RawOptions
func pipelineGeneralGenerator(options RawOptions, globalVars *pipelineGlobalOption) RawOptions {
	if globalVars.AppSpec != nil {
		globalVars.AppSpec.updatePiplineOption(options)
	}
	// update image related config
	newOption := make(RawOptions)
	newOption["pipeline"] = options
	newOption["scm"] = RawOptions(globalVars.Repo.Encode())
	return newOption
}

// pipelineArgocdAppGenerator generate argocdApp options from RawOptions
func pipelineArgocdAppGenerator(options RawOptions, globalVars *pipelineGlobalOption) RawOptions {
	// config app default options
	if _, exist := options["app"]; !exist {
		options["app"] = RawOptions{
			"name":      globalVars.AppName,
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
	if source, sourceExist := options["source"]; sourceExist {
		sourceMap := source.(RawOptions)
		if _, repoURLExist := sourceMap["repoURL"]; !repoURLExist {
			sourceMap["repoURL"] = globalVars.Repo.GetCloneURL()
			sourceMap["repoBranch"] = globalVars.Repo.Branch
		}
		options["source"] = sourceMap
	} else {
		options["source"] = RawOptions{
			"valuefile":  "values.yaml",
			"path":       fmt.Sprintf("helm/%s", globalVars.AppName),
			"repoURL":    string(globalVars.Repo.GetCloneURL()),
			"repoBranch": globalVars.Repo.Branch,
		}
	}
	// config imageRepo default options
	if _, imageRepoExist := options["imageRepo"]; !imageRepoExist {
		if len(globalVars.ImageRepo) > 0 {
			options["imageRepo"] = globalVars.ImageRepo
		}
	}
	return options
}

// hasDefaultConfig check whether
func (o *pipelineOption) hasDefaultConfig() bool {
	return o.defaultConfigLocation != ""
}
