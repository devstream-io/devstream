package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager

	Context("States", func() {
		BeforeEach(func() {
			b, err := backend.GetBackend(backend.BackendLocal)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).NotTo(BeNil())

			smgr, err = statemanager.NewManager(b)
			Expect(err).NotTo(HaveOccurred())
			Expect(smgr).NotTo(BeNil())
		})

		It("Should get the state right", func() {
			key := "name_githubactions"
			stateA := statemanager.State{
				Name:     "name",
				Plugin:   configloader.Plugin{Kind: "githubactions", Version: "0.0.2"},
				Options:  map[string]interface{}{"a": "value"},
				Resource: map[string]interface{}{"a": "value"},
			}

			err := smgr.AddState(key, stateA)
			Expect(err).NotTo(HaveOccurred())

			stateB := smgr.GetState(key)
			Expect(&stateA).To(Equal(stateB))

			err = smgr.DeleteState(key)
			Expect(err).NotTo(HaveOccurred())

			stateC := smgr.GetState(key)
			Expect(stateC).To(BeZero())
		})

		AfterEach(func() {
			err := os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
