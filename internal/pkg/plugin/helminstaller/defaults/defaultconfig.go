package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

var DefaultOptionsMap = make(map[string]*helm.Options)
var StatusGetterFuncMap = make(map[string]plugininstaller.StatusGetterOperation)
