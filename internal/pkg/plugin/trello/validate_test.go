package trello

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("setDefault method", func() {
	var rawOptions configmanager.RawOptions
	When("board is not exist", func() {
		BeforeEach(func() {
			rawOptions = configmanager.RawOptions{
				"scm": configmanager.RawOptions{
					"name":    "test",
					"owner":   "test_user",
					"scmType": "github",
				},
			}
		})
		It("should set default value", func() {
			data, err := setDefault(rawOptions)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data["board"]).Should(Equal(map[string]any{
				"name":        "test_user/test",
				"description": "Description is managed by DevStream, please don't modify. test_user/test",
			}))
			Expect(data["scm"]).Should(Equal(map[string]any{
				"owner":   "test_user",
				"name":    "test",
				"scmType": "github",
			}))
			Expect(data["ci"]).ShouldNot(BeEmpty())
		})
	})
})

var _ = Describe("validate method", func() {
	var rawOptions configmanager.RawOptions
	When("board is not exist", func() {
		BeforeEach(func() {
			rawOptions = configmanager.RawOptions{
				"scm": configmanager.RawOptions{
					"name":    "test",
					"owner":   "test_user",
					"scmType": "github",
				},
			}
		})
		It("should return err", func() {
			_, err := validate(rawOptions)
			Expect(err).Should(HaveOccurred())
		})
	})
})
