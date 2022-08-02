package helm

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/helm"
)

// GetPlugAllStateWrapper will get deploy, ds, statefulset status
func GetPluginAllState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	anFilter := map[string]string{
		helm.GetAnnotationName(): opts.GetReleaseName(),
	}
	labelFilter := map[string]string{}
	return common.GetPluginAllK8sState(opts.GetNamespace(), anFilter, labelFilter)
}
