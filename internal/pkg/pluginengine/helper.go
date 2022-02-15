package pluginengine

import (
	"fmt"
	"os"
	"plugin"
	"time"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/tcnksm/go-input"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func loadPlugin(pluginDir string, tool *configloader.Tool) (DevStreamPlugin, error) {
	mod := fmt.Sprintf("%s/%s", pluginDir, configloader.GetPluginFileName(tool))
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}

	var devStreamPlugin DevStreamPlugin
	symDevStreamPlugin, err := plug.Lookup("DevStreamPlugin")
	if err != nil {
		return nil, err
	}

	devStreamPlugin, ok := symDevStreamPlugin.(DevStreamPlugin)
	if !ok {
		return nil, fmt.Errorf("DevStreamPlugin type error")
	}

	return devStreamPlugin, nil
}

func execute(p *planmanager.Plan) map[string]error {
	errorsMap := make(map[string]error)

	log.Info("Start executing the plan.")
	log.Infof("Changes count: %d.", len(p.Changes))

	for i, c := range p.Changes {
		log.Separatorf("Processing progress: %d/%d.", i+1, len(p.Changes))
		log.Infof("Processing: %s -> %s ...", c.Tool.Name, c.ActionName)

		var succeeded bool
		var err error

		switch c.ActionName {
		case statemanager.ActionInstall:
			if _, err = Create(c.Tool); err == nil {
				succeeded = true
			}
		case statemanager.ActionReinstall:
			if _, err = Update(c.Tool); err == nil {
				succeeded = true
			}
		case statemanager.ActionUninstall:
			succeeded, err = Delete(c.Tool)
		}

		if err != nil {
			key := fmt.Sprintf("%s-%s", c.Tool.Name, c.ActionName)
			errorsMap[key] = err
		}

		c.Result = &planmanager.ChangeResult{
			Succeeded: succeeded,
			Error:     err,
			Time:      time.Now().Format(time.RFC3339),
		}

		err = p.HandleResult(c)
		if err != nil {
			errorsMap["handle-result"] = err
		}
	}

	return errorsMap
}

func readUserInput() string {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Continue? [y/n]"
	userInput, err := ui.Ask(query, &input.Options{
		Required: true,
		Default:  "n",
		Loop:     true,
		ValidateFunc: func(s string) error {
			if s != "y" && s != "n" {
				return fmt.Errorf("input must be y or n")
			}
			return nil
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return userInput
}
