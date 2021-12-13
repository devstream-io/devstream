package plan

import (
	"testing"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func TestMakePlan(t *testing.T) {
	cfg := &config.Config{
		Tools: []config.Tool{
			{
				Name:    "tool_a",
				Version: "v0.0.1",
				Options: map[string]interface{}{"key": "value"},
			},
		},
	}

	b, err := backend.GetBackend("local")
	if err != nil {
		t.Fatal("failed to get backend.")
	}
	smgr := statemanager.NewManager(b)

	smgr.AddState(statemanager.NewState("tool_b", "v0.0.1", nil, statemanager.StatusRunning, &statemanager.Operation{
		Action:   "",
		Time:     "",
		Metadata: map[string]interface{}{"key": "value"},
	}))

	// tool_a should be installed and tool_b should be uninstalled.
	plan := NewPlan(smgr, cfg)

	if len(plan.Changes) != 2 {
		t.Errorf("plan length error")
	}

	for _, change := range plan.Changes {
		if (change.Tool.Name == "tool_a" && change.ActionName == statemanager.ActionInstall) ||
			(change.Tool.Name == "tool_b" && change.ActionName == statemanager.ActionUninstall) {
			continue
		}
		t.Errorf("plan item error")
	}
}
