package template

var readGoNameTpl = "read.go"
var readGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var readGoContentTpl = `package [[ .Name | format ]]

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
		},
		GetStatusOperation: func(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
			// TODO(dtm): Add your GetStatusOperation here.
			return nil, nil
		},
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(options)
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
