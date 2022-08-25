package k8s

import (
	"context"
	"errors"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) ListDeploymentsWithLabel(namespace string, labelFilter map[string]string) ([]appsv1.Deployment, error) {
	dpList, err := c.clientset.AppsV1().Deployments(namespace).List(context.TODO(), generateLabelFilterOption(labelFilter))
	if err != nil {
		return nil, err
	}
	return dpList.Items, nil
}

func (c *Client) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	return c.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) CreateDeployment(namespace string, deployment *appsv1.Deployment) error {
	if _, err := c.clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{}); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The Deployment %s is already exists.", deployment.Name)
	}
	log.Debugf("The Deployment %s has been created.", deployment.Name)
	return nil
}

// Wait for deployment to be ready after creating
func (c *Client) WaitForDeploymentReady(retry int, namespace, deployName string) error {
	deployRunning := false
	for i := 0; i < retry; i++ {
		var dp *appsv1.Deployment
		dp, err := c.GetDeployment(namespace, deployName)
		if err != nil {
			return err
		}

		if isDeploymentReady(dp) {
			log.Infof("The deployment %s is ready.", dp.Name)
			deployRunning = true
			break
		}
		time.Sleep(5 * time.Second)
		log.Debugf("Retry check deployment status %v times", i)
	}

	if !deployRunning {
		return errors.New("create deployment failed")
	}
	return nil
}

func (c *Client) DeleteDeployment(namespace, deployName string) error {
	return c.clientset.AppsV1().Deployments(namespace).
		Delete(context.TODO(), deployName, metav1.DeleteOptions{})
}

func (c *Client) ListDaemonsetsWithLabel(namespace string, labeFilter map[string]string) ([]appsv1.DaemonSet, error) {
	dsList, err := c.clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), generateLabelFilterOption(labeFilter))
	if err != nil {
		return nil, err
	}
	return dsList.Items, nil
}

func (c *Client) GetDaemonset(namespace, name string) (*appsv1.DaemonSet, error) {
	return c.clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) ListStatefulsetsWithLabel(namespace string, labelFilter map[string]string) ([]appsv1.StatefulSet, error) {
	ssList, err := c.clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), generateLabelFilterOption(labelFilter))
	if err != nil {
		return nil, err
	}
	return ssList.Items, nil
}

func (c *Client) GetStatefulset(namespace, name string) (*appsv1.StatefulSet, error) {
	return c.clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func generateLabelFilterOption(labelFilter map[string]string) metav1.ListOptions {
	labelSelector := metav1.LabelSelector{MatchLabels: labelFilter}
	options := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	return options
}

func isDeploymentReady(deployment *appsv1.Deployment) bool {
	return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
}

func isDaemonsetReady(daemonset *appsv1.DaemonSet) bool {
	return daemonset.Status.NumberReady == daemonset.Status.DesiredNumberScheduled
}

func isStatefulsetReady(statefulset *appsv1.StatefulSet) bool {
	return statefulset.Status.ReadyReplicas == *statefulset.Spec.Replicas
}
