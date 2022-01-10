package planmanager_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Planmanager", func() {
	var (
		smgr statemanager.Manager
	)

	BeforeEach(func() {
		defer GinkgoRecover()

		b, err := backend.GetBackend("local")
		Expect(err).NotTo(HaveOccurred())

		smgr = statemanager.NewManager(b)
		_, _ = GinkgoWriter.Write([]byte("new a statemanager"))

		DeferCleanup(func() {
			os.Remove(local.DefaultStateFile)
		})
	})

	It("should be 'one install'", func() {
		name := "tool_a"
		version := "v0.0.1"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, version)},
		}
		plan := planmanager.NewPlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(1))
		c := plan.Changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionInstall))
	})

	It("should be 'two install'", func() {
		name1, name2 := "tool_a", "tool_b"
		version1, version2 := "v0.0.1", "v0.0.2"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name1, version1), *getTool(name2, version2)},
		}
		plan := planmanager.NewPlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(2))

		c1 := plan.Changes[0]
		Expect(c1.Tool.Name).To(Equal(name1))
		Expect(c1.Tool.Version).To(Equal(version1))
		Expect(c1.ActionName).To(Equal(statemanager.ActionInstall))

		c2 := plan.Changes[1]
		Expect(c2.Tool.Name).To(Equal(name2))
		Expect(c2.Tool.Version).To(Equal(version2))
		Expect(c2.ActionName).To(Equal(statemanager.ActionInstall))
	})

	It("should be 1 uninstall when `dtm delete` is triggered against a config with 1 tool and a successful state", func() {
		name := "tool_a"
		version := "v0.0.1"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, version)},
		}
		smgr.AddState(createState(name, version))
		plan := planmanager.NewDeletePlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(1))
		c := plan.Changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionUninstall))
	})
})

func getTool(name, version string) *configloader.Tool {
	return &configloader.Tool{
		Name:    name,
		Version: version,
		Options: map[string]interface{}{"key": "value"},
	}
}

func createState(name, version string) *statemanager.State {
	return &statemanager.State{
		Name:         name,
		Version:      version,
		Dependencies: make([]string, 0),
		Status:       statemanager.StatusInstalled,
		LastOperation: &statemanager.Operation{
			Action:   statemanager.ActionInstall,
			Time:     time.Now().Format(time.RFC3339),
			Metadata: map[string]interface{}{},
		}}
}
