package template

var deleteGoNameTpl = "delete.go"
var deleteGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var deleteGoContentTpl = `package [[ .Name | format ]]

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			// TODO(dtm): Add your ExecuteOperations here.
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	return true, nil
}
`

var delete_go_mustExistFlag = true

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       deleteGoNameTpl,
		DirTpl:        deleteGoDirTpl,
		ContentTpl:    deleteGoContentTpl,
		MustExistFlag: delete_go_mustExistFlag,
	})
}
