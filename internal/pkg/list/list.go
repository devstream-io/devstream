package list

import (
	"github.com/devstream-io/devstream/internal/pkg/list/plugins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// list is the version of DevStream.
// Assign the value when building with the -X parameter. Example:
// -X github.com/devstream-io/devstream/internal/pkg/list.PluginsName=${PLUGINS_NAME}
// See the Makefile for more info.

var PluginsName string

type Action string

const (
	ActionListPlugin Action = "plugins"
)

var ActionSet = map[Action]struct{}{
	ActionListPlugin: {},
}

func IsValideAction(action Action) bool {
	_, ok := ActionSet[action]
	return ok
}

func ExecuteAction(action Action) error {
	switch action {
	case ActionListPlugin:
		log.Debugf("Action: %s.", ActionListPlugin)
		return plugins.List(PluginsName)
	default:
		panic("This should be never happen!")
	}
}
