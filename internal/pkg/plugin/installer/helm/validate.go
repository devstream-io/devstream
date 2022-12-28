package helm

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the dtm-core.
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	var structErrs = validator.CheckStructError(opts)

	if opts.Chart.ChartPath == "" && (opts.Repo.Name == "" || opts.Repo.URL == "" || opts.Chart.ChartName == "") {
		log.Debugf("Repo.Name: %s, Repo.URL: %s, Chart.ChartName: %s", opts.Repo.Name, opts.Repo.URL, opts.Chart.ChartName)
		err := fmt.Errorf("if chartPath == \"\", then the repo.Name & repo.URL & chart.chartName must be set")
		structErrs = append(structErrs, err)
	}

	return options, structErrs.Combine()
}

// SetDefaultConfig will update options empty values base on import options
func SetDefaultConfig(defaultConfig *Options) installer.MutableOperation {
	return func(options configmanager.RawOptions) (configmanager.RawOptions, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}
		opts.FillDefaultValue(defaultConfig)
		return types.EncodeStruct(opts)
	}
}
