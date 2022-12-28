package template

var validateGoNameTpl = "validate.go"
var validateGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var validateGoContentTpl = `package [[ .Name | format ]]

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	if errs := validator.CheckStructError(opts).Combine(); len(errs) != 0 {
		return nil, errs.Combine()
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
