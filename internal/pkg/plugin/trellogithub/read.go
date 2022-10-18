package trellogithub

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	var ws = []*github.Workflow{trelloWorkflow}
	retMap, err := tg.VerifyWorkflows(ws)
	if err != nil {
		return nil, err
	}

	errFlag := false
	for name, err := range retMap {
		if err != nil {
			errFlag = true
			log.Errorf("The workflow/file %s got some error: %s", name, err)
		}
		log.Infof("The workflow/file %s is ok", name)
	}
	if errFlag {
		return nil, nil
	}

	return tg.buildReadStatus()
}
