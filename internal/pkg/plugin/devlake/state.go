package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func getStaticState(opts plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	res := make(map[string]interface{})
	res["deployments"] = make(map[string]interface{})
	res["services"] = make(map[string]interface{})
	for _, d := range devLakeDeployments {
		res["deployments"].(map[string]interface{})[d] = true
		res["services"].(map[string]interface{})[d] = true
	}
	return res, nil
}

func getDynamicState(opts plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	labelFilter := map[string]string{
		"app": "devlake",
	}
	return common.GetPluginAllK8sState(defaultNamespace, map[string]string{}, labelFilter)
}
