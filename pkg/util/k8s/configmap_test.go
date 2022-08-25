package k8s

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("configmap methods", func() {
	var (
		client                   *Client
		configmapName, namespace string
		testConfigMap            []runtime.Object
		labels, data             map[string]string
	)
	BeforeEach(func() {
		configmapName = "test_configmap"
		namespace = "test"
		labels = map[string]string{
			"usage": "test",
		}
		data = map[string]string{
			"field": "test",
		}
		client = &Client{}
		client.clientset = fake.NewSimpleClientset()
	})
	Context("ApplyConfigMap method", func() {
		When("configmap exist", func() {
			BeforeEach(func() {
				testConfigMap = []runtime.Object{&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configmapName,
						Namespace: namespace,
						Labels:    labels,
					},
					Data: data,
				}}
				client.clientset = fake.NewSimpleClientset(testConfigMap...)
			})
			It("should udpate configmap", func() {
				currectConfigMap, err := client.clientset.CoreV1().ConfigMaps(
					namespace).Get(context.Background(), configmapName, metav1.GetOptions{})
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(currectConfigMap.Data).Should(Equal(data))
				Expect(currectConfigMap.ObjectMeta.Labels).Should(Equal(labels))
			})
		})
	})
	Context("GetConfigMap method", func() {
		When("configmap not exist", func() {
			BeforeEach(func() {
				configmapName = "not_exist_config"
			})
			It("should return not found error", func() {
				_, err := client.GetConfigMap(namespace, configmapName)
				Expect(err).Error().Should(HaveOccurred())
				Expect(errors.IsNotFound(err)).Should(BeTrue())
			})
		})
	})
})
