package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("DecodePlugin func", func() {
	var (
		plugData *mockStruct
		rawData  configmanager.RawOptions
	)
	When("decoder is not valid", func() {
		BeforeEach(func() {
			plugData = nil
		})
		It("should return error", func() {
			err := util.DecodePlugin(rawData, plugData)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("create plugin decoder failed"))
		})
	})
	When("options is not valid", func() {
		BeforeEach(func() {
			plugData = new(mockStruct)
			rawData = map[string]any{
				"scm": map[string]any{"key": "not_exist"},
			}
		})
		It("should return error", func() {
			err := util.DecodePlugin(rawData, plugData)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("decode plugin option failed"))
		})
	})

	When("all params are valid", func() {
		BeforeEach(func() {
			plugData = new(mockStruct)
			rawData = map[string]any{
				"scm": map[string]any{
					"url": "github.com/test/test_repo",
				},
			}
		})

		It("should set default value", func() {
			err := util.DecodePlugin(rawData, plugData)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(plugData).Should(Equal(&mockStruct{
				Scm: &git.RepoInfo{
					Owner:    "test",
					Repo:     "test_repo",
					Branch:   "main",
					RepoType: "github",
					CloneURL: "github.com/test/test_repo",
					NeedAuth: false,
				},
			}))
		})
	})
})
