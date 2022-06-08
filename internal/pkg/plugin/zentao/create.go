package zentao

import (
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Retry times for check zentao deployment status, currently this means 5 seconds * 120 times = 10 minutes
const retryTimes int = 120

// Used to rollback when errors occur during creation
var installStep []string = []string{"pv", "pvc", "app"}

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	if err := createByClientAPI(&opts); err != nil {
		return nil, err
	}

	// TODO(southNorth): Integration with other devops tools
	// Currently just store zentao's application status: "running" or "stopped"
	return map[string]interface{}{"running": true}, nil
}

// Create zentao application by go client, maybe other ways will be implemented later
func createByClientAPI(opts *Options) error {

	// 1. Create k8s clientset
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}
	// 2. Create zentao ns
	if err := DealWithNsWhenInstall(kubeClient, opts); err != nil {
		return err
	}

	// Used to indicate steps when errors occur during creation
	step := 0

	// 6. Delete zentao namespace when meet error
	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := DealWithErrWhenInstall(kubeClient, opts, step); err != nil {
			log.Errorf("Failed to deal with creation exit: %s.", err)
		}
		log.Debugf("Deal with creation exit when encounter errors succeeded.")
	}()

	// 3. Create zentao PV
	//    PV will not be recreate in Update
	if retErr = CreatePersistentVolume(kubeClient, opts); retErr != nil {
		return retErr
	}

	// 4. Create zentao PVC
	//    PVC will not be recreate in `Update`
	step++
	if retErr = CreatePersistentVolumeClaim(kubeClient, opts); retErr != nil {
		return retErr
	}

	// 5. Create zentao application
	//    Deploy and service will be recreate in `Update`
	step++
	if retErr = CreateZentaoAPP(kubeClient, opts); retErr != nil {
		return retErr
	}

	return nil
}

// Create zentao application
func CreateZentaoAPP(kubeClient *k8s.Client, opts *Options) error {

	// Create zentao service deployment
	if err := CreateDeployment(kubeClient, opts); err != nil {
		return err
	}

	// Create zentao service
	if err := CreateService(kubeClient, opts); err != nil {
		return err
	}

	// Wait for deployment to be ready
	deployRunning := false
	for i := 0; i < retryTimes; i++ {
		var dp *appsv1.Deployment
		dp, err := kubeClient.GetDeployment(opts.Namespace, opts.Deployment.Name)
		if err != nil {
			return err
		}

		if kubeClient.IsDeploymentReady(dp) {
			log.Infof("The deployment %s is ready.", dp.Name)
			deployRunning = true
			break
		}
		time.Sleep(5 * time.Second)
		log.Debugf("Retry check deployment status %v times", i)
	}

	if !deployRunning {
		return errors.New("create zentao deployment failed")
	}

	return nil
}
