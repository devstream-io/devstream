package template

var readGoNameTpl = "read.go"
var readGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var readGoContentTpl = `package [[ .Name | format ]]

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
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

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    readGoNameTpl,
		DirTpl:     readGoDirTpl,
		ContentTpl: readGoContentTpl,
	})
}
