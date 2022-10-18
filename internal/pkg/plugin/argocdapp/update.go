package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Update(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return Create(options)
}
