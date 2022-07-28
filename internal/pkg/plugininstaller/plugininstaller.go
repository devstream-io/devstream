package plugininstaller

import "github.com/devstream-io/devstream/pkg/util/log"

type RawOptions map[string]interface{}

type (
	// MutableOperation will changes options if it is needed
	MutableOperation func(options RawOptions) (RawOptions, error)
	// BaseOperation only reads options and executes operation
	BaseOperation func(options RawOptions) error
	// StatusOperation only reads options and executes operation
	StatusOperation func(options RawOptions) (map[string]interface{}, error)
)

type Installer interface {
	Execute(options RawOptions) (map[string]interface{}, error)
}

// TODO(daniel-hutao): refactor all caller to use NewInstaller() instead of call Runner.
func NewInstaller(preExecOps []MutableOperation, execOps, termiOps []BaseOperation, getStatusOps StatusOperation) Installer {
	return &Runner{
		PreExecuteOperations: preExecOps,
		ExecuteOperations:    execOps,
		TerminateOperations:  termiOps,
		GetStatusOperation:   getStatusOps,
	}
}

// Runner is the basic type of Installer, It organize func to run in order
type Runner struct {
	PreExecuteOperations []MutableOperation
	ExecuteOperations    []BaseOperation
	TerminateOperations  []BaseOperation
	GetStatusOperation   StatusOperation
}

func (runner *Runner) Execute(options RawOptions) (map[string]interface{}, error) {
	var err error
	// 1. Run PreExecuteOperations first, these func can change options
	log.Debugf("Start Execute PreInstall Operations...")
	for _, preInstallOperation := range runner.PreExecuteOperations {
		options, err = preInstallOperation(options)
		if err != nil {
			return nil, err
		}
	}
	// 2. register terminate function if encounter in install
	var installError error
	defer func() {
		if installError == nil {
			return
		}
		log.Debugf("Start to execute terminating operations...")
		for _, terminateOperation := range runner.TerminateOperations {
			err := terminateOperation(options)
			if err != nil {
				log.Errorf("Failed to deal with namespace: %s.", err)
			}
		}
	}()

	log.Debugf("Start to execute install operations...")
	// 3. Run ExecuteOperations in order, these func can't change options
	for _, installOperation := range runner.ExecuteOperations {
		installError = installOperation(options)
		if installError != nil {
			return nil, installError
		}
	}
	// 4. Get Status for this execution step
	var status map[string]interface{}
	if runner.GetStatusOperation != nil {
		log.Debugf("Start to execute getting status operations...")
		status, err = runner.GetStatusOperation(options)
	}
	return status, err
}
