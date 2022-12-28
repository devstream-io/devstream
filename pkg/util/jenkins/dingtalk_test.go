package jenkins

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

var _ = Describe("jenkins dingtalk methods", func() {
	var (
		s              *ghttp.Server
		j              JenkinsAPI
		err            error
		dingtalkConfig *dingtalk.BotConfig
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
		dingtalkConfig = &dingtalk.BotConfig{
			RobotConfigs: []dingtalk.BotInfoConfig{
				{
					ID:      "test",
					Name:    "test",
					Webhook: "test",
				},
			},
		}
	})
	Context("ApplyDingTalkBot method", func() {
		When("apply success", func() {
			BeforeEach(func() {
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/manage/dingtalk/configure"),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)
			})
			It("should work normal", func() {
				err := j.ApplyDingTalkBot(*dingtalkConfig)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
		When("apply failed", func() {
			BeforeEach(func() {
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/manage/dingtalk/configure"),
						ghttp.RespondWith(http.StatusBadGateway, "bad gateway"),
					),
				)
			})
			It("should return error", func() {
				err := j.ApplyDingTalkBot(*dingtalkConfig)
				Expect(err).Error().Should(HaveOccurred())
			})
		})

	})
	AfterEach(func() {
		s.Close()
	})
})
