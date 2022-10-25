package jenkins

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("jenkinsOption struct", func() {
	Context("newClient method", func() {
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
		When("env password is setted", func() {
			var (
				existEnvPassword, testPassword string
			)
			BeforeEach(func() {
				testPassword = "test"
				existEnvPassword = os.Getenv(jenkinsPasswordEnvKey)
				os.Setenv(jenkinsPasswordEnvKey, testPassword)
			})
			It("should work normal", func() {
				auth, err := opt.getBasicAuth()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(auth).ShouldNot(BeNil())
				Expect(auth.Username).Should(Equal(user))
				Expect(auth.Password).Should(Equal(testPassword))
			})
			AfterEach(func() {
				if existEnvPassword != "" {
					os.Setenv(jenkinsPasswordEnvKey, existEnvPassword)
				} else {
					os.Unsetenv(jenkinsPasswordEnvKey)
				}
			})
		})
	})
})
