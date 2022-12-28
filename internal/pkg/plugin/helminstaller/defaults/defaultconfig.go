package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
)

var (
	DefaultOptionsMap   = make(map[string]*helm.Options)
	StatusGetterFuncMap = make(map[string]installer.StatusGetterOperation)
)
