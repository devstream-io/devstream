package plugininstaller

import (
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type RawOptions map[string]interface{}

type (
	// MutableOperation can be used to change options if it is needed
	MutableOperation func(options RawOptions) (RawOptions, error)
	// BaseOperation reads options and executes operation
	BaseOperation func(options RawOptions) error
	// StateOperation reads options and executes operation, then returns the state map
	StateOperation func(options RawOptions) (statemanager.ResourceState, error)
)

type (
	PreExecuteOperations []MutableOperation
	ExecuteOperations    []BaseOperation
	TerminateOperations  []BaseOperation
	GetStateOperation    StateOperation
)

type Installer interface {
	Execute(options RawOptions) (map[string]interface{}, error)
}

// Operator knows all the operations and can execute them in order
type Operator struct {
	PreExecuteOperations PreExecuteOperations
	ExecuteOperations    ExecuteOperations
	TerminateOperations  TerminateOperations
	GetStateOperation    GetStateOperation
}

// Execute will sequentially execute all operations in Operator
func (o *Operator) Execute(options RawOptions) (map[string]interface{}, error) {
	var err error
	// 1. Execute PreExecuteOperations. It may changes the options.
	log.Debugf("Start to execute PreExecuteOperations...")
	for _, preOps := range o.PreExecuteOperations {
		options, err = preOps(options)
		if err != nil {
			return nil, err
		}
	}

	// 2. Register defer func so that in case ExecuteOperations fails, it can execute TerminateOperations
	var execErr error
	defer func() {
		if execErr == nil {
			return
		}
		log.Debugf("Start to execute TerminateOperations...")
		for _, terminateOperation := range o.TerminateOperations {
			err := terminateOperation(options)
			if err != nil {
				log.Errorf("Failed to execute TerminateOperations: %s.", err)
			}
		}
	}()

	// 3. Execute ExecuteOperations in order. It won't change the options.
	log.Debugf("Start to execute ExecuteOperations...")
	for _, execOps := range o.ExecuteOperations {
		execErr = execOps(options)
		if execErr != nil {
			return nil, execErr
		}
	}

	// 4. Execute GetStateOperation.
	var state map[string]interface{}
	if o.GetStateOperation != nil {
		log.Debugf("Start to execute GetStateOperation...")
		state, err = o.GetStateOperation(options)
		if err != nil {
			return nil, err
		}
	}
	return state, err
}
