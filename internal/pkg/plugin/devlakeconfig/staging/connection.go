package staging

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/devlakeconfig/staging/common"
)

// BaseConnection FIXME ...
type BaseConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	common.Model
}

// BasicAuth FIXME ...
type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password" encrypt:"yes"`
}

// AccessToken FIXME ...
type AccessToken struct {
	Token string `mapstructure:"token" validate:"required" json:"token" encrypt:"yes"`
}

// AppKey FIXME ...
type AppKey struct {
	AppId     string `mapstructure:"app_id" validate:"required" json:"appId"`
	SecretKey string `mapstructure:"secret_key" validate:"required" json:"secretKey" encrypt:"yes"`
}

// RestConnection FIXME ...
type RestConnection struct {
	BaseConnection   `mapstructure:",squash"`
	Endpoint         string `mapstructure:"endpoint" validate:"required" json:"endpoint"`
	Proxy            string `mapstructure:"proxy" json:"proxy"`
	RateLimitPerHour int    `comment:"api request rate limit per hour" json:"rateLimitPerHour"`
}
