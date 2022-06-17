package tekton

import (
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var DefaultDeploymentList = []string{
	"tekton-pipelines-controller",
	"tekton-pipelines-webhook",
}

func GetStaticState() *helm.InstanceState {
	retState := &helm.InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}

func GetDynamicState(opts *Options) (*helm.InstanceState, error) {
	namespace := getOrDefaultNamespace(opts.Chart.Namespace)
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		return nil, err
	}

	retState := &helm.InstanceState{}
	for _, dp := range dps {
		dpAn := dp.GetAnnotations()
		dpName := dp.GetName()
		helmNameAn, exist := dpAn[GetAnnotationName()]
		if !exist || helmNameAn != opts.Chart.ReleaseName {
			log.Infof("Found unknown deployment: %s.", dp.GetName())
		}
		ready := kubeClient.IsDeploymentReady(&dp)
		retState.Workflows.AddDeployment(dpName, ready)
		log.Debugf("The deployment %s is %t.", dp.GetName(), ready)
	}
	return retState, err
}

func getOrDefaultNamespace(namespace string) string {
	if namespace == "" {
		//tekon is default namespace
		namespace = "tekon"
	}
	return namespace
}
