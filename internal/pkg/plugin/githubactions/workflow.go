package githubactions

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/plugin/githubactions/golang"
	"github.com/merico-dev/stream/internal/pkg/plugin/githubactions/nodejs"
	"github.com/merico-dev/stream/internal/pkg/plugin/githubactions/python"
)

const (
	defaultCommitMessage = "builder by DevStream"
	BuilderYmlPr         = "pr-builder.yml"
	BuilderYmlMaster     = "master-builder.yml"
	BuilderYmlPipeline   = "pipeline.yml"
)

var go117 = &Language{
	Name:    "go",
	Version: "1.17",
}

var python3 = &Language{
	Name:    "python",
	Version: "3",
}

var nodejs9 = &Language{
	Name:    "nodejs",
	Version: "9",
}

var defaultWorkflows = workflows{
	go117.String(): {
		{defaultCommitMessage, BuilderYmlPipeline, golang.PipelineBuilder},
	},
	python3.String(): {
		{defaultCommitMessage, BuilderYmlPr, python.PrBuilder},
		{defaultCommitMessage, BuilderYmlMaster, python.MasterBuilder},
	},
	nodejs9.String(): {
		{defaultCommitMessage, BuilderYmlPr, nodejs.PrBuilder},
		{defaultCommitMessage, BuilderYmlMaster, nodejs.MasterBuilder},
	},
}

// TODO(daniel-hutao): Options should keep as same as other plugins named Param
// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string
	Repo     string
	Branch   string
	Language *Language
	Jobs     *Jobs
}

// Language is the struct containing details of a programming language specified in the GitHub Actions Workflow.
type Language struct {
	Name    string
	Version string
}

// Workflow is the struct for a GitHub Actions Workflow.
type Workflow struct {
	commitMessage    string
	workflowFileName string
	workflowContent  string
}

type LanguageString string

type workflows map[LanguageString][]*Workflow

func (l *Language) Validate() []error {
	retErrors := make([]error, 0)

	if l.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("name is empty"))
	}

	return retErrors
}

func (l *Language) String() LanguageString {
	return LanguageString(fmt.Sprintf("%s-%s", l.Name, l.Version))
}

func (ws *workflows) GetWorkflowByNameVersionTypeString(nvtStr LanguageString) []*Workflow {
	workflowList, exist := (*ws)[nvtStr]
	if exist {
		return workflowList
	}
	return nil
}
