package docker

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func ComposeUp(options plugininstaller.RawOptions) error {
	if err := op.ComposeUp(); err != nil {
		return fmt.Errorf("failed to run containers: %s", err)
	}

	return nil
}

func ComposeDown(options plugininstaller.RawOptions) error {
	if err := op.ComposeDown(); err != nil {
		return fmt.Errorf("failed to down containers: %s", err)
	}

	return nil
}

func ComposeState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	state, err := op.ComposeState()
	if err != nil {
		return nil, fmt.Errorf("failed to get containers state: %s", err)
	}

	return state, nil
}
