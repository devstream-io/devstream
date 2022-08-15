package template

var createGoNameTpl = "create.go"
var createGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var createGoContentTpl = `package [[ .Name | format ]]

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

    // TODO(dtm): Add your logic here.

    return nil, nil
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
