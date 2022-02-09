package statemanager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Statemanager", func() {
	var smgr statemanager.Manager

	Context("States", func() {
		BeforeEach(func() {
			b, err := backend.GetBackend("local")
			Expect(err).NotTo(HaveOccurred())
			Expect(b).NotTo(BeNil())

			smgr = statemanager.NewManager(b)
			Expect(smgr).NotTo(BeNil())
		})

		It("Should get the state right", func() {
			stateA := statemanager.NewState("prod", configloader.Plugin{Kind: "argocd", Version: "v0.0.1"}, nil, nil)
			Expect(stateA).NotTo(BeNil())

			smgr.AddState(stateA)

			stateB := smgr.GetState("prod_argocd")
			Expect(stateA).To(Equal(stateB))

			smgr.DeleteState("prod_argocd")
			stateC := smgr.GetState("prod_argocd")
			Expect(stateC).To(BeNil())
		})

		It("Should Read/Write well", func() {
			// write
			stateA := statemanager.NewState("prod", configloader.Plugin{Kind: "argocd", Version: "v0.0.1"}, []string{}, map[string]interface{}{})
			smgr.AddState(stateA)
			err := smgr.Write(smgr.GetStatesMap().Format())
			Expect(err).NotTo(HaveOccurred())

			// read
			data, err := smgr.Read()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).NotTo(BeZero())

			tmpMap := make(map[string]*statemanager.State)
			err = yaml.Unmarshal(data, tmpMap)
			Expect(err).NotTo(HaveOccurred())

			statesMap := statemanager.NewStatesMap()
			for k, v := range tmpMap {
				statesMap.Store(k, v)
			}

			stateB, ok := statesMap.Load("prod_argocd")
			Expect(ok).To(BeTrue())
			Expect(*stateB.(*statemanager.State)).To(Equal(*stateA))
		})

		AfterEach(func() {
			err := os.RemoveAll(local.DefaultStateFile)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
