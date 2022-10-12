package devlakeconfig

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugin/devlakeconfig/staging"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

// Options is the struct for configurations of the devlake-config plugin.
type Options struct {
	DevLakeAddr string          `mapstructure:"devlakeAddr" validate:"url"`
	Plugins     []DevLakePlugin `mapstructure:"plugins" validate:"required"`
}

// NewOptions create options by raw options
func NewOptions(options plugininstaller.RawOptions) (Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}

type DevLakePlugin struct {
	Name        string       `mapstructure:"name" validate:"required"`
	Connections []Connection `mapstructure:"connections"`
}

// TODO(daniel-hutao): uncomment the code below after DevLake fix the umimportable issue:
//
//github.com/devstream-io/devstream/internal/pkg/plugin/devlakeconfig imports
//        github.com/apache/incubator-devlake/plugins/helper tested by
//        github.com/apache/incubator-devlake/plugins/helper.test imports
//        github.com/apache/incubator-devlake/mocks: module github.com/apache/incubator-devlake@latest found (v0.14.0), but does not contain package github.com/apache/incubator-devlake/mocks
//
//type Connection struct {
//	helper.RestConnection `mapstructure:",squash"`
//	helper.BasicAuth      `mapstructure:",squash"`
//	helper.AccessToken    `mapstructure:",squash"`
//	helper.AppKey         `mapstructure:",squash"`
//}

type Connection struct {
	staging.RestConnection `mapstructure:",squash"`
	Authx                  Auth `mapstructure:"auth" validate:"required"`
	Auth                   `mapstructure:",squash"`
}

type Auth struct {
	staging.BasicAuth   `mapstructure:",squash"`
	staging.AccessToken `mapstructure:",squash"`
	staging.AppKey      `mapstructure:",squash"`
}

func (o *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(o, &options); err != nil {
		return nil, err
	}
	return options, nil
}
