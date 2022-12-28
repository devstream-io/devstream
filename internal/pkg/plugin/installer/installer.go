package installer

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type (
	// MutableOperation can be used to change options if it is needed
	MutableOperation func(options configmanager.RawOptions) (configmanager.RawOptions, error)
	// BaseOperation reads options and executes operation
	BaseOperation func(options configmanager.RawOptions) error
	// StatusGetterOperation reads options and executes operation, then returns the status map
	StatusGetterOperation func(options configmanager.RawOptions) (statemanager.ResourceStatus, error)
)

type (
	PreExecuteOperations []MutableOperation
	ExecuteOperations    []BaseOperation
	TerminateOperations  []BaseOperation
)

type Installer interface {
	Execute(options configmanager.RawOptions) (statemanager.ResourceStatus, error)
}

// Operator knows all the operations and can execute them in order
type Operator struct {
	PreExecuteOperations PreExecuteOperations
	ExecuteOperations    ExecuteOperations
	TerminateOperations  TerminateOperations
	GetStatusOperation   StatusGetterOperation
}

// Execute will sequentially execute all operations in Operator
func (o *Operator) Execute(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
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

	// 4. Execute StatusGetterOperation.
	var state map[string]interface{}
	if o.GetStatusOperation != nil {
		log.Debugf("Start to execute StatusGetterOperation...")
		state, err = o.GetStatusOperation(options)
		if err != nil {
			return nil, err
		}
	}
	return state, err
}
