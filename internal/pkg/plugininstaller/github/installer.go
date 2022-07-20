package github

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func ProcessAction(action string) plugininstaller.BaseOperation {
	return func(options plugininstaller.RawOptions) error {
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
			switch action {
			case "create":
				err = ghClient.AddWorkflow(w, opts.Branch)
			case "delete":
				err = ghClient.DeleteWorkflow(w, opts.Branch)
			case "update":
				err = ghClient.DeleteWorkflow(w, opts.Branch)
				if err != nil {
					err = ghClient.AddWorkflow(w, opts.Branch)
				}
			default:
				err = fmt.Errorf("This github Action not support: %s", action)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}
