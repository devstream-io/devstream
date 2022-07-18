package plugininstaller

import "github.com/devstream-io/devstream/pkg/util/log"

type RawOptions map[string]interface{}

// MutableOperation will change options if need
type MutableOperation func(options RawOptions) (RawOptions, error)

// BaseOperation only read options and execute operation
type BaseOperation func(options RawOptions) error

// StatusOperation only read options and execute operation
type StatusOperation func(options RawOptions) (map[string]interface{}, error)

// Runner is the basic type of plugininstaller, It organize func to run in order
type Runner struct {
	PreExecuteOperations []MutableOperation
	ExecuteOperations    []BaseOperation
	TermateOperations    []BaseOperation
	GetStatusOperation   StatusOperation
}

func (runner *Runner) Execute(options RawOptions) (map[string]interface{}, error) {
	var err error
	// 1. Run PreExecuteOperations first, these func can change options
	for _, preInstallOperation := range runner.PreExecuteOperations {
		options, err = preInstallOperation(options)
		if err != nil {
			return nil, err
		}
	}
	// 2. register termate function if encounter in install
	var installError error
	defer func() {
		if installError == nil {
			return
		}
		for _, termateOperation := range runner.TermateOperations {
			err := termateOperation(options)
			if err != nil {
				log.Errorf("Failed to deal with namespace: %s.", err)
			}
		}
		log.Debugf("Deal with termate func succeeded.")
	}()

	// 3. Run ExecuteOperations in order, these func can't change options
	for _, installOperation := range runner.ExecuteOperations {
		installError = installOperation(options)
		if err != nil {
			break
		}
	}
	// 4. If encounter err, return error
	if installError != nil {
		return nil, err
	}
	// 5. Get Status for this execute
	var status map[string]interface{}
	if runner.GetStatusOperation != nil {
		status, err = runner.GetStatusOperation(options)
	}
	return status, err
}
