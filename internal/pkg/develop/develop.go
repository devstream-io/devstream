package develop

import (
	"github.com/merico-dev/stream/internal/pkg/develop/plugin"
	"github.com/merico-dev/stream/pkg/util/log"
)

type Action string

const (
	ActionCreatePlugin Action = "create-plugin"
)

var ActionSet = map[Action]struct{}{
	ActionCreatePlugin: {},
}

func IsValideAction(action Action) bool {
	_, ok := ActionSet[action]
	return ok
}

func BranchAction(action Action) error {
	switch action {
	case ActionCreatePlugin:
		log.Debugf("Action: %s", ActionCreatePlugin)
		return plugin.Create()
	default:
		panic("This should be never happen!")
	}
}
