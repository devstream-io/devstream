package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
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

		It("Should get the state list", func() {
			key := statemanager.StateKey("a_githubactions")
			stateA := statemanager.State{
				Name:     "a",
				Plugin:   "githubactions",
				Options:  map[string]interface{}{"a": "value"},
				Resource: map[string]interface{}{"a": "value"},
			}
			err = smgr.AddState(key, stateA)
			Expect(err).NotTo(HaveOccurred())

			key = statemanager.StateKey("b_githubactions")
			stateB := statemanager.State{
				Name:     "b",
				Plugin:   "githubactions",
				Options:  map[string]interface{}{"b": "value"},
				Resource: map[string]interface{}{"b": "value"},
			}
			err = smgr.AddState(key, stateB)
			Expect(err).NotTo(HaveOccurred())

			stateList := smgr.GetStatesMap().ToList()
			Expect(stateList).To(Equal([]statemanager.State{stateA, stateB}))
		})

		AfterEach(func() {
			err = os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
