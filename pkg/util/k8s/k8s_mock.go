package k8s

import (
	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type MockClient struct {
	GetSecretError error
	GetSecretValue map[string]string
}

func (m *MockClient) GetSecret(namespace, name string) (map[string]string, error) {
	if m.GetSecretError != nil {
		return nil, m.GetSecretError
	}
	return m.GetSecretValue, nil
}
func (m *MockClient) ApplySecret(name, namespace string, data map[string][]byte, labels map[string]string) (*corev1.Secret, error) {
	return nil, nil
}
func (m *MockClient) CreateService(namespace string, service *corev1.Service) error {
	return nil
}
func (m *MockClient) DeleteService(namespace, serviceName string) error {
	return nil
}
func (m *MockClient) GetService(namespace, name string) (*corev1.Service, error) {
	return nil, nil
}
func (m *MockClient) CreatePersistentVolume(option *PVOption) error {
	return nil
}
func (m *MockClient) DeletePersistentVolume(pvName string) error {
	return nil
}
func (m *MockClient) CreatePersistentVolumeClaim(opt *PVCOption) error {
	return nil
}
func (m *MockClient) DeletePersistentVolumeClaim(namespace, pvcName string) error {
	return nil
}
func (m *MockClient) GetResourceStatus(nameSpace string, anFilter, labelFilter map[string]string) (*AllResourceStatus, error) {
	return nil, nil
}
func (m *MockClient) ListDeploymentsWithLabel(namespace string, labelFilter map[string]string) ([]appsv1.Deployment, error) {
	return nil, nil
}
func (m *MockClient) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	return nil, nil
}
func (m *MockClient) CreateDeployment(namespace string, deployment *appsv1.Deployment) error {
	return nil
}
func (m *MockClient) WaitForDeploymentReady(retry int, namespace, deployName string) error {
	return nil
}
func (m *MockClient) DeleteDeployment(namespace, deployName string) error {
	return nil
}

func (m *MockClient) ListDaemonsetsWithLabel(namespace string, labeFilter map[string]string) ([]appsv1.DaemonSet, error) {
	return nil, nil
}

func (m *MockClient) GetStatefulset(namespace, name string) (*appsv1.StatefulSet, error) {
	return nil, nil
}

func (m *MockClient) UpsertNameSpace(nameSpace string) error {
	return nil
}

func (m *MockClient) GetNamespace(namespace string) (*corev1.Namespace, error) {
	return nil, nil
}

func (m *MockClient) IsDevstreamNS(namespace string) (bool, error) {
	return false, nil
}

func (m *MockClient) CreateNamespace(namespace string) error {
	return nil
}

func (m *MockClient) DeleteNamespace(namespace string) error {
	return nil
}

func (m *MockClient) IsNamespaceExists(namespace string) (bool, error) {
	return false, nil
}

func (m *MockClient) ApplyConfigMap(name, namespace string, data, labels map[string]string) (*corev1.ConfigMap, error) {
	return nil, nil
}

func (m *MockClient) GetConfigMap(name, namespace string) (*corev1.ConfigMap, error) {
	return nil, nil
}

func (m *MockClient) ListArgocdApplications(namespace string) ([]argocdv1alpha1.Application, error) {
	return nil, nil
}

func (m *MockClient) GetArgocdApplication(namespace, name string) (*argocdv1alpha1.Application, error) {
	return nil, nil
}

func (m *MockClient) IsArgocdApplicationReady(application *argocdv1alpha1.Application) bool {
	return false
}

func (m *MockClient) DescribeArgocdApp(app *argocdv1alpha1.Application) map[string]interface{} {
	return nil
}
