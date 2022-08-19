package configmanager_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("Core Config", func() {
	Context("Validate method", func() {
		It("should return error if state config is empty", func() {
			coreConfig := configmanager.CoreConfig{State: nil}
			err := coreConfig.ValidateAndDefault()
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("state config is empty"))
		})

		It("should return error if backend not exist", func() {
			coreConfig := configmanager.CoreConfig{
				State: &configmanager.State{
					Backend: "not_exist",
				},
			}
			err := coreConfig.ValidateAndDefault()
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(
				ContainSubstring("backend type error"),
			)
		})

		It("should return error s3 option not config", func() {
			coreConfig := configmanager.CoreConfig{
				State: &configmanager.State{Backend: "s3"},
			}
			err := coreConfig.ValidateAndDefault()
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(
				ContainSubstring("state s3 Bucket is empty"),
			)
		})

		It("should return true if config local valid", func() {
			coreConfig := configmanager.CoreConfig{
				State: &configmanager.State{Backend: "local"},
			}
			err := coreConfig.ValidateAndDefault()
			Expect(err).Error().ShouldNot(HaveOccurred())
		})

		It("should return true if config s3 valid", func() {
			coreConfig := configmanager.CoreConfig{
				State: &configmanager.State{
					Backend: "s3",
					Options: configmanager.StateConfigOptions{
						Bucket:    "test_bucket",
						Region:    "test_region",
						Key:       "test_key",
						StateFile: "test_file",
					},
				},
			}
			err := coreConfig.ValidateAndDefault()
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	Context("ParseToolFilePath method", func() {
		It("should return nil if varFile is empty", func() {
			coreConfig := configmanager.CoreConfig{State: nil}
			err := coreConfig.ParseToolFilePath()
			Expect(err).Should(BeNil())
		})

		It("should return error if varFile not exist", func() {
			notExistFile := "not_exist_file"
			coreConfig := configmanager.CoreConfig{State: nil, VarFile: notExistFile}
			err := coreConfig.ParseToolFilePath()
			Expect(err).Error().Should(HaveOccurred())
		})

	})
})
