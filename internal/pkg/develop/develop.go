package develop

import (
	"github.com/devstream-io/devstream/internal/pkg/develop/plugin"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Action string

const (
	ActionCreatePlugin   Action = "create-plugin"
	ActionValidatePlugin Action = "validate-plugin"
)

var ActionSet = map[Action]struct{}{
	ActionCreatePlugin:   {},
	ActionValidatePlugin: {},
}

func IsValideAction(action Action) bool {
	_, ok := ActionSet[action]
	return ok
}

func CreatePlugin() error {
	log.Debugf("Start create plugin")
	return plugin.Create()
}

func ValidatePlugin() error {
	log.Debugf("Start validate plugin")
	return plugin.Validate()
}
