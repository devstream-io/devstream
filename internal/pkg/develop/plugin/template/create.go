package template

var createGoNameTpl = "create.go"
var createGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var createGoContentTpl = `package [[ .Name | format ]]

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			// TODO(dtm): Add your ExecuteOperations here.
		},
		TerminateOperations: plugininstaller.TerminateOperations{
			// TODO(dtm): Add your TerminateOperations here.
		},
		GetStateOperation: func(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
			// TODO(dtm): Add your GetStateOperation here.
			return nil, nil
		},
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
`
var create_go_mustExistFlag = true

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       createGoNameTpl,
		DirTpl:        createGoDirTpl,
		ContentTpl:    createGoContentTpl,
		MustExistFlag: create_go_mustExistFlag,
	})
}
