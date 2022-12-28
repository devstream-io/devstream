package docker

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func ComposeUp(options configmanager.RawOptions) error {
	if err := op.ComposeUp(); err != nil {
		return fmt.Errorf("failed to run containers: %s", err)
	}

	return nil
}

func ComposeDown(options configmanager.RawOptions) error {
	if err := op.ComposeDown(); err != nil {
		return fmt.Errorf("failed to down containers: %s", err)
	}

	return nil
}

func ComposeStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	state, err := op.ComposeState()
	if err != nil {
		return nil, fmt.Errorf("failed to get containers state: %s", err)
	}

	return state, nil
}
