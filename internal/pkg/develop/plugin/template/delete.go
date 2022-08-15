package template

var deleteGoNameTpl = "delete.go"
var deleteGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var deleteGoContentTpl = `package [[ .Name | format ]]

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	// TODO(dtm): Add your logic here.

    return false, nil
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
