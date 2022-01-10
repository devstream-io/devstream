package pluginengine

import (
	"fmt"
	"log"
	"plugin"
	"time"

	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func loadPlugin(tool *configloader.Tool) (DevStreamPlugin, error) {
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return nil, fmt.Errorf("plugin-dir is \"\"")
	}

	mod := fmt.Sprintf("%s/%s_%s.so", pluginDir, tool.Name, tool.Version)
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
		return nil, err
	}

	return devStreamPlugin, nil
}

func execute(p *planmanager.Plan) map[string]error {
	errorsMap := make(map[string]error)

	log.Printf("Changes count: %d.", len(p.Changes))

	for i, c := range p.Changes {
		log.Printf("Processing progress: %d/%d.", i+1, len(p.Changes))
		log.Printf("Processing: %s -> %s ...", c.Tool.Name, c.ActionName)

		var succeeded bool
		var err error

		switch c.ActionName {
		case statemanager.ActionInstall:
			succeeded, err = Install(c.Tool)
		case statemanager.ActionReinstall:
			succeeded, err = Reinstall(c.Tool)
		case statemanager.ActionUninstall:
			succeeded, err = Uninstall(c.Tool)
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
