package backend_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("GetBackend", func() {
	When("use local backend", func() {
		It("should return local backend struct", func() {
			state := configmanager.State{Backend: "local"}
			_, err := backend.GetBackend(state)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// TODO: add mock s3 test
	When("use unknown backend", func() {
		It("should return err", func() {
			state := configmanager.State{Backend: "not_exist_plug"}
			_, err := backend.GetBackend(state)
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("s3 config is empty", func() {
		It("should return err of backendOptionErr", func() {
			state := configmanager.State{Backend: "s3"}
			_, err := backend.GetBackend(state)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})
