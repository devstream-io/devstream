package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager
	var err error

	Context("States", func() {
		BeforeEach(func() {
			smgr, err = statemanager.NewManager()
			Expect(err).NotTo(HaveOccurred())
			Expect(smgr).NotTo(BeNil())
		})

		It("Should get the state right", func() {
			key := statemanager.StateKey("name_githubactions")
			stateA := statemanager.State{
				Name:     "name",
				Plugin:   "githubactions",
				Options:  map[string]interface{}{"a": "value"},
				Resource: map[string]interface{}{"a": "value"},
			}

			err = smgr.AddState(key, stateA)
			Expect(err).NotTo(HaveOccurred())

			stateB := smgr.GetState(key)
			Expect(&stateA).To(Equal(stateB))

			err = smgr.DeleteState(key)
			Expect(err).NotTo(HaveOccurred())

			stateC := smgr.GetState(key)
			Expect(stateC).To(BeZero())
		})

		AfterEach(func() {
			err = os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
