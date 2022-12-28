package k8s

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("k8s secret methods", func() {
	var (
		client                                      *Client
		namespace, secretName, secretKey, secretVal string
		secretData                                  map[string][]byte
	)
	BeforeEach(func() {
		client = &Client{}
		namespace = "test"
		secretName = "secret_name"
		client.clientset = fake.NewSimpleClientset()
		secretKey = "testSecret"
		secretVal = "this is a test"
		secretData = map[string][]byte{
			secretKey: []byte(secretVal),
		}
	})
	Context("GetSecret method", func() {
		BeforeEach(func() {
			testResources := []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      secretName,
						Namespace: namespace,
					},
					Data: secretData,
				},
			}
			client.clientset = fake.NewSimpleClientset(testResources...)
		})
		When("return secret", func() {
			It("should return map", func() {
				secret, err := client.GetSecret(namespace, secretName)
				Expect(err).ShouldNot(HaveOccurred())
				val, ok := secret[secretKey]
				Expect(ok).Should(BeTrue())
				Expect(val).Should(Equal(string(secretVal)))
			})
		})
	})
})
