package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager
	var err error

	Context("States", func() {
		BeforeEach(func() {
			stateCfg := configloader.State{
				Backend: "local",
				Options: configloader.StateConfigOptions{
					StateFile: "devstream.state",
				},
			}
			smgr, err = statemanager.NewManager(stateCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(smgr).NotTo(BeNil())
		})

		It("Should get the state right", func() {
			key := statemanager.StateKey("name_githubactions")
			stateA := statemanager.State{
				InstanceID: "name",
				Name:       "githubactions",
				Options:    map[string]interface{}{"a": "value"},
				Resource:   map[string]interface{}{"a": "value"},
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
			// Adding order: A,C,B
			// List order should be: A,B,C
			key := statemanager.StateKey("a_githubactions")
			stateA := statemanager.State{
				InstanceID: "a",
				Name:       "githubactions",
				Options:    map[string]interface{}{"a": "value"},
				Resource:   map[string]interface{}{"a": "value"},
			}
			err = smgr.AddState(key, stateA)
			Expect(err).NotTo(HaveOccurred())

			key = statemanager.StateKey("c_githubactions")
			stateC := statemanager.State{
				InstanceID: "c",
				Name:       "githubactions",
				Options:    map[string]interface{}{"c": "value"},
				Resource:   map[string]interface{}{"c": "value"},
			}
			err = smgr.AddState(key, stateC)
			Expect(err).NotTo(HaveOccurred())

			key = statemanager.StateKey("b_githubactions")
			stateB := statemanager.State{
				InstanceID: "b",
				Name:       "githubactions",
				Options:    map[string]interface{}{"b": "value"},
				Resource:   map[string]interface{}{"b": "value"},
			}
			err = smgr.AddState(key, stateB)
			Expect(err).NotTo(HaveOccurred())

			stateList := smgr.GetStatesMap().ToList()
			Expect(stateList).To(Equal([]statemanager.State{stateA, stateB, stateC}))
		})

		AfterEach(func() {
			err = os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
