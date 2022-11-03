package jenkinspipeline

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/k8s"
)

var _ = Describe("jenkinsOption struct", func() {
	var (
		url, user, namespace string
		opt                  *jenkinsOption
	)
	BeforeEach(func() {
		url = "http://test.exmaple.com"
		user = "test_user"
		namespace = "test_namespace"
		opt = &jenkinsOption{
			URL:       url,
			User:      user,
			Namespace: namespace,
		}
	})
	var (
		existEnvPassword, testPassword string
	)

	BeforeEach(func() {
		existEnvPassword = os.Getenv(jenkinsPasswordEnvKey)
		err := os.Unsetenv(jenkinsPasswordEnvKey)
		Expect(err).Error().ShouldNot(HaveOccurred())
		testPassword = "test"
	})

	Context("getBasicAuth method", func() {
		BeforeEach(func() {
			os.Setenv(jenkinsPasswordEnvKey, testPassword)
		})
		When("env password is setted", func() {
			It("should work normal", func() {
				auth, err := opt.getBasicAuth()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(auth).ShouldNot(BeNil())
				Expect(auth.Username).Should(Equal(user))
				Expect(auth.Password).Should(Equal(testPassword))
			})
		})
	})
	Context("newClient method", func() {
		var (
			s *ghttp.Server
		)
		BeforeEach(func() {
			s = ghttp.NewServer()
			s.RouteToHandler("GET", "/api/json", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "ok")
			})
			s.RouteToHandler("GET", "/crumbIssuer/api/json/api/json", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "ok")
			})
			opt.URL = s.URL()
			os.Setenv(jenkinsPasswordEnvKey, testPassword)
		})

		It("should work normal", func() {
			_, err := opt.newClient()
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			s.Close()
		})
	})
	Context("getAuthFromSecret method", func() {
		var (
			m         *k8s.MockClient
			nameSpace string
		)
		BeforeEach(func() {
			nameSpace = "test_namespace"
		})
		When("GetSecret return error", func() {
			BeforeEach(func() {
				m = &k8s.MockClient{
					GetSecretError: fmt.Errorf("test error"),
				}
			})
			It("should return nil", func() {
				Expect(getAuthFromSecret(m, nameSpace)).Should(BeNil())
			})
		})
		When("defaultUsername not exist in secret", func() {
			BeforeEach(func() {
				m = &k8s.MockClient{
					GetSecretValue: map[string]string{
						defaultAdminSecretUserPassword: "test",
					},
				}
			})
			It("should return nil", func() {
				Expect(getAuthFromSecret(m, nameSpace)).Should(BeNil())
			})
		})
		When("defaultPassword not exist in secret", func() {
			BeforeEach(func() {
				m = &k8s.MockClient{
					GetSecretValue: map[string]string{
						defaultAdminSecretUserName: "test",
					},
				}
			})
			It("should return nil", func() {
				Expect(getAuthFromSecret(m, nameSpace)).Should(BeNil())
			})
		})
		When("GetSecret return error", func() {
			nameSpace = "test_namespace"
			BeforeEach(func() {
				m = &k8s.MockClient{
					GetSecretValue: map[string]string{
						defaultAdminSecretUserName:     "test_user",
						defaultAdminSecretUserPassword: "test_pass",
					},
				}
			})
			It("should return nil", func() {
				auth := getAuthFromSecret(m, nameSpace)
				Expect(auth).ShouldNot(BeNil())
				Expect(auth.Username).Should(Equal("test_user"))
				Expect(auth.Password).Should(Equal("test_pass"))
			})
		})
	})
	AfterEach(func() {
		if existEnvPassword != "" {
			os.Setenv(jenkinsPasswordEnvKey, existEnvPassword)
		}
	})
})
