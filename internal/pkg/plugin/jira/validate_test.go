package jira

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("validate func", func() {
	var (
		rawOpt configmanager.RawOptions
	)
	When("input params is not valid", func() {
		BeforeEach(func() {
			rawOpt = configmanager.RawOptions{}
		})
		It("should return error", func() {
			_, err := validate(rawOpt)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("opts are illegal"))
		})
	})
	When("scm is not github type", func() {
		BeforeEach(func() {
			rawOpt = configmanager.RawOptions{
				"scm": map[string]any{
					"scmType": "gitlab",
					"owner":   "user",
					"name":    "repo",
				},
				"integOptions": map[string]any{
					"baseUrl":    "test.com",
					"userEmail":  "test@test.com",
					"projectKey": "project",
				},
			}
		})
		It("should return error", func() {
			_, err := validate(rawOpt)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("plugin jira-integ only support scm type github for now"))
		})
	})
	When("all params are right", func() {
		BeforeEach(func() {
			rawOpt = configmanager.RawOptions{
				"scm": map[string]any{
					"scmType": "github",
					"owner":   "user",
					"name":    "repo",
				},
				"integOptions": map[string]any{
					"baseUrl":    "test.com",
					"userEmail":  "test@test.com",
					"projectKey": "project",
				},
			}
		})
		It("should return error", func() {
			_, err := validate(rawOpt)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("setDefault func", func() {
	var (
		rawOpt configmanager.RawOptions
	)
	BeforeEach(func() {
		rawOpt = configmanager.RawOptions{
			"scm": map[string]any{
				"scmType": "github",
				"owner":   "user",
				"name":    "repo",
			},
			"integOptions": map[string]any{
				"baseUrl":    "test.com",
				"userEmail":  "test@test.com",
				"projectKey": "project",
			},
		}

	})
	It("should set ci info", func() {
		op, err := setDefault(rawOpt)
		Expect(err).ShouldNot(HaveOccurred())
		ciConfig, exist := op["ci"]
		Expect(exist).Should(BeTrue())
		Expect(ciConfig).ShouldNot(BeNil())
	})
})
