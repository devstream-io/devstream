package jenkins

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Basic Auth", func() {
	var (
		auth     *BasicAuth
		userName string
	)
	Context("CheckNameMatch method", func() {
		When("name is matched", func() {
			BeforeEach(func() {
				auth = &BasicAuth{
					Username: userName,
				}
			})
			It("should return true", func() {
				Expect(auth.CheckNameMatch(userName)).Should(BeTrue())
			})
		})
	})
	Context("usePassWordAuth method", func() {
		When("user name and password are not empty", func() {
			BeforeEach(func() {
				auth = &BasicAuth{
					Username: "test",
					Password: "test",
				}
			})
			It("should return true", func() {
				Expect(auth.usePassWordAuth()).Should(BeTrue())
			})
		})
		When("password is empty", func() {
			BeforeEach(func() {
				auth = &BasicAuth{
					Username: "test",
				}
			})
			It("should return false", func() {
				Expect(auth.usePassWordAuth()).Should(BeFalse())
			})
		})
	})
})

var _ = Describe("jenkins auth methods", func() {
	var (
		s                                                   *ghttp.Server
		j                                                   JenkinsAPI
		credName, credXmlStr, credStatusPath, createReqPath string
		err                                                 error
	)
	BeforeEach(func() {
		s = ghttp.NewServer()
		s.RouteToHandler("GET", "/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		s.RouteToHandler("GET", "/crumbIssuer/api/json/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		opts := &JenkinsConfigOption{
			URL:       s.URL(),
			Namespace: "test",
			BasicAuth: &BasicAuth{
				Username: "test_user",
				Password: "test_password",
			},
		}
		j, err = NewClient(opts)
		Expect(err).ShouldNot(HaveOccurred())
		createReqPath = "/credentials/store/system/domain/_/createCredentials"
		credXmlStr = `<org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl plugin="plain-credentials@139.ved2b_9cf7587b">
<scope>GLOBAL</scope>
<id>sonarqubeTokenCredential</id>
<description>sonarqubeTokenCredential</description>
<secret>
<secret-redacted/>
</secret>
</org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl>`

	})
	Context("CreateGiltabCredential method", func() {
		var token string
		BeforeEach(func() {
			credName = "test_gitlab_cred"
			token = "test_token"
			credStatusPath = fmt.Sprintf("/credentials/store/system/domain/_/credential/%s/config.xml/", credName)
			expectReqBody := fmt.Sprintf(`<com.dabsquared.gitlabjenkins.connection.GitLabApiTokenImpl><id>%s</id><scope>GLOBAL</scope><description>%s</description><apiToken>%s</apiToken></com.dabsquared.gitlabjenkins.connection.GitLabApiTokenImpl>`, credName, credName, token)
			s.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", createReqPath),
					ghttp.VerifyBody([]byte(expectReqBody)),
					ghttp.RespondWith(http.StatusOK, "ok"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", credStatusPath),
					ghttp.RespondWith(http.StatusOK, credXmlStr),
				),
			)
		})
		It("should work normal", func() {
			err := j.CreateGiltabCredential(credName, token)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	Context("CreateSecretCredential method", func() {
		When("cred already exist", func() {
			var secretText string
			BeforeEach(func() {
				credName = "test_secret_cred"
				secretText = "test_secret"
				credStatusPath = fmt.Sprintf("/credentials/store/system/domain/_/credential/%s/config.xml", credName)
				expectReqBody := fmt.Sprintf(`<org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl><id>%s</id><scope>GLOBAL</scope><description>%s</description><secret>%s</secret></org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl>`, credName, credName, secretText)
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", createReqPath),
						ghttp.VerifyBody([]byte(expectReqBody)),
						ghttp.RespondWith(http.StatusConflict, "conflict"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", credStatusPath),
						ghttp.VerifyBody([]byte(expectReqBody)),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)
			})
			It("should update cred", func() {
				err := j.CreateSecretCredential(credName, secretText)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("CreatePasswordCredential method", func() {
		var userName, password string
		When("jenkins return error", func() {
			BeforeEach(func() {
				credName = "test_pass_cred"
				userName = "testUser"
				password = "testPass"
				expectReqBody := fmt.Sprintf(`<com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl><id>%s</id><scope>GLOBAL</scope><description>%s</description><username>%s</username><password>%s</password></com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl>`, credName, credName, userName, password)
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", createReqPath),
						ghttp.VerifyBody([]byte(expectReqBody)),
						ghttp.RespondWith(http.StatusBadGateway, "gateway error"),
					),
				)
			})
			It("should return error", func() {
				err := j.CreatePasswordCredential(credName, userName, password)
				Expect(err).Error().Should(HaveOccurred())
			})

		})
	})

	Context("CreateSSHKeyCredential method", func() {
		When("get cred failed", func() {
			var sshKey, userName string
			BeforeEach(func() {
				credName = "test_sshKey_cred"
				sshKey = "test_ssh_key"
				userName = "test_ssh_user"
				credStatusPath = fmt.Sprintf("/credentials/store/system/domain/_/credential/%s/config.xml/", credName)
				expectReqBody := fmt.Sprintf(`<com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey><id>%s</id><scope>GLOBAL</scope><username>%s</username><description>%s</description><privateKeySource class="com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource"><privateKey>%s</privateKey></privateKeySource></com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey>`, credName, userName, credName, sshKey)
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", createReqPath),
						ghttp.VerifyBody([]byte(expectReqBody)),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", credStatusPath),
						ghttp.RespondWith(http.StatusBadGateway, "bad gateway"),
					),
				)
			})
			It("should return error", func() {
				err := j.CreateSSHKeyCredential(credName, userName, sshKey)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	AfterEach(func() {
		s.Close()
	})
})
