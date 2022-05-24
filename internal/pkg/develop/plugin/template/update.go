package template

var update_go_nameTpl = "update.go"
var update_go_dirTpl = "internal/pkg/plugin/{{ .Name | dirFormat }}/"
var update_go_contentTpl = `package {{ .Name | format }}

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)


func Update(options map[string]interface{}) (map[string]interface{}, error) {
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

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    update_go_nameTpl,
		DirTpl:     update_go_dirTpl,
		ContentTpl: update_go_contentTpl,
	})
}
