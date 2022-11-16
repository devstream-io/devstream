package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/helm"
)

var DefaultOptionsMap = make(map[string]*helm.Options)
var StatusGetterFuncMap = make(map[string]plugininstaller.StatusGetterOperation)
