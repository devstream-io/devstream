package step_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("DingtalkStepConfig", func() {
	var (
		c          *step.DingtalkStepConfig
		mockClient *jenkins.MockClient
		name       string
	)
	BeforeEach(func() {
		name = "test"
		c = &step.DingtalkStepConfig{
			Name: name,
		}
	})
	Context("GetJenkinsPlugins method", func() {
		It("should return dingding plugin", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("dingding-notifications"))
		})
	})

	Context("ConfigJenkins method", func() {
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				_, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("ConfigSCM method", func() {
		var (
			scmClient *scm.MockScmClient
			errMsg    string
		)
		When("webhook is not valid", func() {
			BeforeEach(func() {
				errMsg = "step scm dingTalk.webhook is not valid"
				c.Webhook = "test.example.com/gg"
				scmClient = &scm.MockScmClient{}
			})
			It("should return error", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("add token failed", func() {
			BeforeEach(func() {
				errMsg = "add token failed"
				c.Webhook = "test.example.com/gg?token=test"
				scmClient = &scm.MockScmClient{
					AddRepoSecretError: fmt.Errorf(errMsg),
				}
			})
			It("should return error", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("all valid", func() {
			BeforeEach(func() {
				c.Webhook = "test.example.com/gg?token=test"
				scmClient = &scm.MockScmClient{}
			})
			It("should return nil", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})
})
