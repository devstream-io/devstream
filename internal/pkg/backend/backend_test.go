package backend_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("GetBackend", func() {
	When("use local backend", func() {
		It("should return local backend struct", func() {
			state := configmanager.State{Backend: "local"}
			_, err := backend.GetBackend(state)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})

	// TODO: add mock s3 test
	When("use unknown backend", func() {
		It("should return err", func() {
			state := configmanager.State{Backend: "not_exist_plug"}
			_, err := backend.GetBackend(state)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("the backend type < not_exist_plug > is illegal"))
		})
	})
})
