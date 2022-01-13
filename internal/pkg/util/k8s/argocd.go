package k8s

import (
	"context"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/gitops-engine/pkg/health"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ListArgocdApplications(namespace string) ([]argocdv1alpha1.Application, error) {
	appList, err := c.Argocd.ArgoprojV1alpha1().Applications(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return appList.Items, nil
}

func (c *Client) GetArgocdApplication(namespace, name string) (*argocdv1alpha1.Application, error) {
	return c.Argocd.ArgoprojV1alpha1().Applications(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) IsArgocdApplicationReady(application *argocdv1alpha1.Application) bool {
	return application.Status.Health.Status == health.HealthStatusHealthy
}
