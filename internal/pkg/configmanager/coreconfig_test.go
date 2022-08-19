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
			err := coreConfig.Validate()
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("state config is empty"))
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
