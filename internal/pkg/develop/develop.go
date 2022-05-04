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

func ExecuteAction(action Action) error {
	switch action {
	case ActionCreatePlugin:
		log.Debugf("Action: %s.", ActionCreatePlugin)
		return plugin.Create()
	case ActionValidatePlugin:
		log.Debugf("Action: %s.", ActionValidatePlugin)
		return plugin.Validate()
	default:
		panic("This should be never happen!")
	}
}
