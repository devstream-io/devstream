package k8s

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("service methods", func() {
	var (
		client                 *Client
		serviceName, namespace string
	)
	BeforeEach(func() {
		serviceName = "testService"
		namespace = "test"
		client = &Client{}
		client.clientset = fake.NewSimpleClientset()
	})
	Context("GetService method", func() {
		BeforeEach(func() {
			testNameSpace := []runtime.Object{
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceName,
						Namespace: namespace,
					},
				},
			}
			client.clientset = fake.NewSimpleClientset(testNameSpace...)
		})

		It("should return service", func() {
			service, err := client.GetService(namespace, serviceName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(service.Name).Should(Equal(serviceName))
		})
	})
	Context("DeleteService method", func() {
		When("service not exist", func() {
			BeforeEach(func() {
				serviceName = "not_exist"
			})
			It("should return error", func() {
				err := client.DeleteService(namespace, serviceName)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("service exist", func() {
			BeforeEach(func() {
				testNameSpace := []runtime.Object{
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:      serviceName,
							Namespace: namespace,
						},
					},
				}
				client.clientset = fake.NewSimpleClientset(testNameSpace...)
			})

			It("should delete service", func() {
				err := client.DeleteService(namespace, serviceName)
				Expect(err).Error().ShouldNot(HaveOccurred())
				_, err = client.clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("CreateService method", func() {
		It("should work normal", func() {
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceName,
					Namespace: namespace,
				},
			}
			err := client.CreateService(namespace, service)
			Expect(err).Error().ShouldNot(HaveOccurred())
			_, err = client.clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})
