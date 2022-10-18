package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager
	var err error
	var testKey statemanager.StateKey
	var testState statemanager.State

	BeforeEach(func() {
		stateCfg := configmanager.State{
			Backend: "local",
			Options: configmanager.StateConfigOptions{
				StateFile: "devstream.state",
			},
		}
		smgr, err = statemanager.NewManager(stateCfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(smgr).NotTo(BeNil())
		testKey = statemanager.StateKey("name_githubactions")
		testState = statemanager.State{
			InstanceID:     "name",
			Name:           "githubactions",
			Options:        configmanager.RawOptions{"a": "value"},
			ResourceStatus: statemanager.ResourceStatus{"a": "value"},
		}
	})

	Describe("Manager", func() {
		It("Should add state right", func() {
			err = smgr.AddState(testKey, testState)
			Expect(err).NotTo(HaveOccurred())

			stateB := smgr.GetState(testKey)
			Expect(&testState).To(Equal(stateB))

			err = smgr.DeleteState(testKey)
			Expect(err).NotTo(HaveOccurred())

			stateC := smgr.GetState(testKey)
			Expect(stateC).To(BeZero())
		})

		It("Should get the state list", func() {
			// Adding order: A,C,B
			// List order should be: A,B,C
			key := statemanager.StateKey("a_githubactions")
			stateA := statemanager.State{
				InstanceID:     "a",
				Name:           "githubactions",
				Options:        map[string]interface{}{"a": "value"},
				ResourceStatus: map[string]interface{}{"a": "value"},
			}
			err = smgr.AddState(key, stateA)
			Expect(err).NotTo(HaveOccurred())

			key = statemanager.StateKey("c_githubactions")
			stateC := statemanager.State{
				InstanceID:     "c",
				Name:           "githubactions",
				Options:        configmanager.RawOptions{"c": "value"},
				ResourceStatus: statemanager.ResourceStatus{"c": "value"},
			}
			err = smgr.AddState(key, stateC)
			Expect(err).NotTo(HaveOccurred())

			key = statemanager.StateKey("b_githubactions")
			stateB := statemanager.State{
				InstanceID:     "b",
				Name:           "githubactions",
				Options:        configmanager.RawOptions{"b": "value"},
				ResourceStatus: statemanager.ResourceStatus{"b": "value"},
			}
			err = smgr.AddState(key, stateB)
			Expect(err).NotTo(HaveOccurred())

			stateList := smgr.GetStatesMap().ToList()
			Expect(stateList).To(Equal([]statemanager.State{stateA, stateB, stateC}))
		})
	})

	AfterEach(func() {
		err = os.RemoveAll(local.DefaultStateFile)
		Expect(err).NotTo(HaveOccurred())
	})
})
