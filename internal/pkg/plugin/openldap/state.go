package openldap

import "github.com/merico-dev/stream/pkg/util/helm"

var DefaultDeploymentList = []string{
	"openldap-ltb-passwd",
	"openldap-phpldapadmin",
}

func GetStaticState() *helm.InstanceState {
	retState := &helm.InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}
