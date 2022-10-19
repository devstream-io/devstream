package devlakeconfig

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/devlakeconfig/staging"
)

// Options is the struct for configurations of the devlake-config plugin.
type Options struct {
	DevLakeAddr string          `mapstructure:"devlakeAddr" validate:"url"`
	Plugins     []DevLakePlugin `mapstructure:"plugins" validate:"required"`
}

// NewOptions create options by raw options
func NewOptions(options configmanager.RawOptions) (Options, error) {
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
	// a connection format in dtm config:
	// - name: ""
	//   endpoint: ""
	//   proxy: ""
	//   rateLimitPerHour: 0
	//   auth:
	//     username: "changeme"
	//     password: "changeme"
	Auth Auth `mapstructure:"auth" validate:"required"`
	// a connection format in DevLake api:
	// {
	//   "name": ""
	//   "endpoint": ""
	//   "proxy": ""
	//   "rateLimitPerHour": 0
	//   "username": "changeme"
	//   "password": "changeme"
	// }
	InlineAuth `mapstructure:",squash"`
}

type Auth struct {
	staging.BasicAuth   `mapstructure:",squash"`
	staging.AccessToken `mapstructure:",squash"`
	staging.AppKey      `mapstructure:",squash"`
}

type InlineAuth Auth

func (o *Options) Encode() (configmanager.RawOptions, error) {
	var options configmanager.RawOptions
	if err := mapstructure.Decode(o, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func RenderAuthConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	for _, p := range opts.Plugins {
		for _, c := range p.Connections {
			c.Token = c.Auth.Token
			c.Username = c.Auth.Username
			c.Password = c.Auth.Password
			c.AppId = c.Auth.AppId
			c.SecretKey = c.Auth.SecretKey
		}
	}

	return opts.Encode()
}
