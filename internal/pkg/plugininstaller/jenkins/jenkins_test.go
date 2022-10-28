package jenkins

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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
	AfterEach(func() {
		if existEnvPassword != "" {
			os.Setenv(jenkinsPasswordEnvKey, existEnvPassword)
		}
	})
})
