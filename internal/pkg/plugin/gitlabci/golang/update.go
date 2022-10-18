package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return Create(options)
}
