package plugininstaller

import "github.com/devstream-io/devstream/pkg/util/log"

type RawOptions map[string]interface{}

type (
	// MutableOperation can be used to change options if it is needed
	MutableOperation func(options RawOptions) (RawOptions, error)
	// BaseOperation reads options and executes operation
	BaseOperation func(options RawOptions) error
	// StateOperation reads options and executes operation, then returns the state map
	StateOperation func(options RawOptions) (map[string]interface{}, error)
)

type Installer interface {
	Execute(options RawOptions) (map[string]interface{}, error)
}

// TODO(daniel-hutao): refactor all caller to use NewInstaller() instead of call Operator.
func NewInstaller(preExecOps []MutableOperation, execOps, termiOps []BaseOperation, getStatusOps StateOperation) Installer {
	return &Operator{
		PreExecuteOperations: preExecOps,
		ExecuteOperations:    execOps,
		TerminateOperations:  termiOps,
		GetStateOperation:    getStatusOps,
	}
}

// Operator knows all the operations and can execute them in order
type Operator struct {
	PreExecuteOperations []MutableOperation
	ExecuteOperations    []BaseOperation
	TerminateOperations  []BaseOperation
	GetStateOperation    StateOperation
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
	}
	return state, err
}
