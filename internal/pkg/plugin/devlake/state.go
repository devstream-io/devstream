package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/util"
)

func getStaticState(opts plugininstaller.RawOptions) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	res["deployments"] = make(map[string]interface{})
	res["services"] = make(map[string]interface{})
	for _, d := range devLakeDeployments {
		res["deployments"].(map[string]interface{})[d] = true
		res["services"].(map[string]interface{})[d] = true
	}
	return res, nil
}

func getDynamicState(opts plugininstaller.RawOptions) (map[string]interface{}, error) {
	return util.ReadDepAndServiceState(defaultNamespace, devLakeDeployments)
}
