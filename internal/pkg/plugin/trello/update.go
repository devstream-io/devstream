package trello

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

// Create creates Tello board and lists(todo/doing/done).
func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return Create(options)
}
