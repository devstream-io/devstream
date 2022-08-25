package k8s

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("namespace methods", func() {
	var (
		client         *Client
		namespace      string
		devstreamLabel map[string]string
	)
	BeforeEach(func() {
		namespace = "test"
		client = &Client{}
		client.clientset = fake.NewSimpleClientset()
		devstreamLabel = map[string]string{
			"created_by": "DevStream",
		}
	})
	Context("UpsertNameSpace method", func() {
		When("namespace not exist", func() {
			BeforeEach(func() {
				namespace = "not_exist"
			})
			It("should create namespace", func() {
				err := client.UpsertNameSpace(namespace)
				Expect(err).Error().ShouldNot(HaveOccurred())
				namespaceData, err := client.clientset.CoreV1().Namespaces().Get(
					context.TODO(), namespace, metav1.GetOptions{},
				)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(namespaceData.Name).Should(Equal(namespace))
			})
		})
		When("namespace exist", func() {
			BeforeEach(func() {
				testNameSpace := []runtime.Object{&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace,
					},
				}}
				client.clientset = fake.NewSimpleClientset(testNameSpace...)
			})
			It("should return nil error", func() {
				err := client.UpsertNameSpace(namespace)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})
	Context("IsDevstreamNS method", func() {
		BeforeEach(func() {
			testNameSpace := []runtime.Object{
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace,
					},
				},
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "test1",
						Labels: devstreamLabel,
					},
				},
				&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "test2",
						Labels: devstreamLabel,
					},
				},
			}
			client.clientset = fake.NewSimpleClientset(testNameSpace...)
		})

		It("should check is devstream namespace", func() {
			result, err := client.IsDevstreamNS("test1")
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(result).Should(BeTrue())
			result, err = client.IsDevstreamNS(namespace)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(result).Should(BeFalse())
		})
	})
	Context("CreateNamespace method", func() {
		It("should create namespace with devstream label", func() {
			err := client.CreateNamespace(namespace)
			Expect(err).Error().ShouldNot(HaveOccurred())
			namespaceData, err := client.clientset.CoreV1().Namespaces().Get(
				context.TODO(), namespace, metav1.GetOptions{},
			)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(namespaceData.Name).Should(Equal(namespace))
			labels := namespaceData.Labels
			createUser, ok := labels["created_by"]
			Expect(ok).Should(BeTrue())
			Expect(createUser).Should(Equal("DevStream"))
		})
	})
	Context("DeleteNamespace method", func() {
		When("namespace is default", func() {
			BeforeEach(func() {
				namespace = "default"
			})
			It("should reutrn err", func() {
				err := client.DeleteNamespace(namespace)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("namespace is valid", func() {
			BeforeEach(func() {
				testNameSpace := []runtime.Object{&corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace,
					},
				}}
				client.clientset = fake.NewSimpleClientset(testNameSpace...)
			})
			It("should return nil error", func() {
				err := client.DeleteNamespace(namespace)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})
	Context("IsNamespaceExists method", func() {
		BeforeEach(func() {
			testNameSpace := []runtime.Object{&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}}
			client.clientset = fake.NewSimpleClientset(testNameSpace...)
		})
		When("namespace not exist", func() {
			BeforeEach(func() {
				namespace = "not_exist"
			})
			It("should return false with no error", func() {
				exist, err := client.IsNamespaceExists(namespace)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(exist).Should(BeFalse())
			})
		})
		When("namespace exist", func() {
			It("should return true with no error", func() {
				exist, err := client.IsNamespaceExists(namespace)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(exist).Should(BeTrue())
			})
		})
	})
})
