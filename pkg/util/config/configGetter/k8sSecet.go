package configGetter

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/devstream-io/devstream/pkg/util/k8s"
)

type k8sSecretGetter struct {
	key                   string
	namespace, secretName string
	K8sClient             kubernetes.Interface
}

func NewK8sSecretGetter(key, namespace, secretName string) *k8sSecretGetter {
	return &k8sSecretGetter{
		namespace:  namespace,
		secretName: secretName,
		key:        key,
	}
}

func (g *k8sSecretGetter) Get() string {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return ""
	}
	secret, err := k8sClient.GetSecret(g.namespace, g.secretName)
	if err != nil {
		// log.Warnf("failed to get secret <%s/%s>: %v", g.namespace, g.secretName, err)
		return ""
	}
	value, ok := secret[g.key]
	if !ok {
		return ""
	}
	return value
}

func (g *k8sSecretGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in k8s secret <%s/%s>", g.key, g.namespace, g.secretName)
}
