package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/plan"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var installCMD = &cobra.Command{
	Use:   "install",
	Short: "Install tools defined in config file",
	Long:  `Install tools defined in config file`,
	Run:   installCMDFunc,
}

func installCMDFunc(cmd *cobra.Command, args []string) {
	conf := config.LoadConf(configFile)
	// use default local backend for now.
	b, err := backend.GetBackend("local")
	if err != nil {
		log.Fatal(err)
	}
	smgr := statemanager.NewManager(b)

	// this channel will be consume quickly, so its size is no need too big.
	execResultChan := make(chan plan.ExecutionResult, 1)

	p := plan.MakePlan(smgr, conf)
	if p.Changes == nil || len(p.Changes) == 0 {
		log.Println("it is nothing to do here")
		return
	}
	go p.Execute(execResultChan)
	handleResult(smgr, execResultChan)

	log.Printf("===")
}

func handleResult(smgr statemanager.Manager, execResultChan <-chan plan.ExecutionResult) {
	var state *statemanager.State
	// recordState record a single tool's state and return if this state is Succeeded.
	recordState := func(smgr statemanager.Manager, change *plan.Change) bool {
		if change.Result.Error != nil {
			state = statemanager.NewState(
				change.Tool.Name,
				change.Tool.Version,
				[]string{},
				statemanager.StatusFailed,
				&statemanager.Operation{
					Action:   change.ActionName,
					Time:     change.Result.Time,
					Metadata: change.Tool.Options,
				},
			)
			smgr.AddState(state)
			log.Printf("=== plugin %s installation failed ===", change.Tool.Name)
			return false
		}
		state = statemanager.NewState(
			change.Tool.Name,
			change.Tool.Version,
			[]string{},
			statemanager.StatusRunning,
			&statemanager.Operation{
				Action:   change.ActionName,
				Time:     change.Result.Time,
				Metadata: change.Tool.Options,
			},
		)
		smgr.AddState(state)
		log.Printf("=== plugin %s installation done ===", change.Tool.Name)
		return true
	}

	var hasToolFailed bool
	for result := range execResultChan {
		for _, change := range result {
			hasToolFailed = !recordState(smgr, change)
		}
	}

	err := smgr.Write(smgr.GetStates().Format())
	if err != nil {
		log.Fatal(err)
	}

	if hasToolFailed {
		log.Printf("=== some plugin installation failed ===")
	}
}
