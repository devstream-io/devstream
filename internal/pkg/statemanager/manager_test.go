package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager

	Context("States", func() {
		BeforeEach(func() {
			b, err := backend.GetBackend(backend.BackendLocal)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).NotTo(BeNil())

			smgr = statemanager.NewManager(b)
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

			smgr.AddState(key, stateA)

			stateB := smgr.GetState(key)
			Expect(&stateA).To(Equal(stateB))

			smgr.DeleteState(key)
			stateC := smgr.GetState(key)
			Expect(stateC).To(BeZero())
		})

		It("Should Read/Write well", func() {
			// write
			key := "state-a"
			stateA := statemanager.State{
				Name:     "name",
				Plugin:   configloader.Plugin{Kind: "githubactions", Version: "0.0.2"},
				Options:  map[string]interface{}{"a": "value"},
				Resource: map[string]interface{}{"a": "value"},
			}
			smgr.AddState(key, stateA)
			err := smgr.Write(smgr.GetStatesMap().Format())
			Expect(err).NotTo(HaveOccurred())

			// read
			data, err := smgr.Read()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).NotTo(BeZero())

			tmpMap := make(map[string]statemanager.State)
			err = yaml.Unmarshal(data, tmpMap)
			Expect(err).NotTo(HaveOccurred())
			log.Infof("tmpMap: %v", tmpMap)

			statesMap := statemanager.NewStatesMap()
			for k, v := range tmpMap {
				statesMap.Store(k, v)
			}

			stateB, ok := statesMap.Load(key)
			Expect(ok).To(BeTrue())
			Expect(stateB.(statemanager.State)).To(Equal(stateA))
		})

		AfterEach(func() {
			err := os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
