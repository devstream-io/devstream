package template

var validateGoNameTpl = "validate.go"
var validateGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var validateGoContentTpl = `package [[ .Name | format ]]

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	errs := validator.Struct(opts)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}
	return options, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    validateGoNameTpl,
		DirTpl:     validateGoDirTpl,
		ContentTpl: validateGoContentTpl,
	})
}
