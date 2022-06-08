package zentao

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	if err := deleteByClientAPI(&opts); err != nil {
		return false, err
	}

	return true, nil
}

// Delete everything created by zentao plugins: service,deployment,pvc,pv,namespace
func deleteByClientAPI(opts *Options) error {
	// 1. Create k8s clientset
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// 2. Delete zentao application
	log.Debug("Prepare to delete zentao application.")
	if err := DeleteZentaoAPP(kubeClient, opts); err != nil {
		return err
	}

	// 3. Delete zentao PVC
	log.Debug("Prepare to delete zentao PVC.")
	if err = DeletePersistentVolumeClaim(kubeClient, opts); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	// 4. Delete zentao PV
	log.Debug("Prepare to delete zentao PV.")
	if err = DeletePersistentVolume(kubeClient, opts); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	// 5. Delete zentao namespace only when namespace is controlled by dtm
	exist, err := IsDevstreamNSExists(kubeClient, opts.Namespace)
	if err != nil {
		log.Debugf("Failed to check whether namespace: %s exists.", opts.Namespace)
		return err
	}

	if exist {
		log.Debug("Prepare to delete zentao namespace.")
		if err = DeleteNamespace(kubeClient, opts); err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		}
	}

	return nil
}

func DeleteZentaoAPP(kubeClient *k8s.Client, opts *Options) error {
	// Delete zentao service
	if err := DeleteService(kubeClient, opts); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	// Delete zentao deployment
	if err := DeleteDeployment(kubeClient, opts); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	return nil
}
