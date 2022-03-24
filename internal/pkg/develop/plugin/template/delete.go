package template

var delete_go_nameTpl = "delete.go"
var delete_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var delete_go_contentTpl = `package {{ .Name }}

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
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

    return false, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    delete_go_nameTpl,
		DirTpl:     delete_go_dirTpl,
		ContentTpl: delete_go_contentTpl,
	})
}
