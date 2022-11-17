package s3

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/types"
)

var _ = Describe("validate", func() {
	It("should return error s3 option not config", func() {
		err := validate("", "", "")
		Expect(err).Error().Should(HaveOccurred())
		Expect(types.IsBackendOptionErr(err)).To(BeTrue())
	})

	It("should return true if config s3 valid", func() {
		bucket := "test_bucket"
		region := "test_region"
		key := "test_key"

		err := validate(bucket, region, key)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
})
