package github

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	ActionCreate action = "create"
	ActionDelete action = "delete"
	ActionUpdate action = "update"
)

type action string

// ProcessAction config github
func ProcessAction(act action) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {
		opts, err := NewGithubActionOptions(options)
		if err != nil {
			return err
		}

		log.Debugf("Language is: %s.", opts.GetLanguage())
		ghClient, err := opts.GetGithubClient()
		if err != nil {
			return err
		}

		for _, w := range opts.Workflows {
			switch act {
			case ActionCreate:
				err = ghClient.AddWorkflow(w, opts.Branch)
			case ActionDelete:
				err = ghClient.DeleteWorkflow(w, opts.Branch)
			case ActionUpdate:
				err = ghClient.DeleteWorkflow(w, opts.Branch)
				if err != nil {
					err = ghClient.AddWorkflow(w, opts.Branch)
				}
			default:
				err = fmt.Errorf("This github Action not support: %s", act)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}
