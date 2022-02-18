package planmanager_test

import (
	"os"

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
		name := "a"
		kind := "tool-a"
		version := "v0.0.1"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, kind, version)},
		}
		plan := planmanager.NewPlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(1))
		c := plan.Changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Plugin.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionInstall))
	})

	It("should be 'two install'", func() {
		name1, name2 := "a", "b"
		kind1, kind2 := "tool-a", "too-b"
		version1, version2 := "v0.0.1", "v0.0.2"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name1, kind1, version1), *getTool(name2, kind2, version2)},
		}
		plan := planmanager.NewPlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(2))

		c1 := plan.Changes[0]
		Expect(c1.Tool.Name).To(Equal(name1))
		Expect(c1.Tool.Plugin.Kind).To(Equal(kind1))
		Expect(c1.Tool.Plugin.Version).To(Equal(version1))
		Expect(c1.ActionName).To(Equal(statemanager.ActionInstall))

		c2 := plan.Changes[1]
		Expect(c2.Tool.Name).To(Equal(name2))
		Expect(c2.Tool.Plugin.Kind).To(Equal(kind2))
		Expect(c2.Tool.Plugin.Version).To(Equal(version2))
		Expect(c2.ActionName).To(Equal(statemanager.ActionInstall))
	})

	It("should be 1 uninstall when `dtm delete` is triggered against a config with 1 tool and a successful state", func() {
		name := "a"
		kind := "tool-a"
		version := "v0.0.1"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, kind, version)},
		}
		// TODO(daniel-hutao) wait for refactor
		smgr.AddState("todo", statemanager.State{})
		plan := planmanager.NewDeletePlan(smgr, cfg)

		Expect(len(plan.Changes)).To(Equal(1))
		c := plan.Changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Plugin.Kind).To(Equal(kind))
		Expect(c.Tool.Plugin.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionUninstall))
	})
})

func getTool(name, kind, version string) *configloader.Tool {
	return &configloader.Tool{
		Name:    name,
		Plugin:  configloader.Plugin{Kind: kind, Version: version},
		Options: map[string]interface{}{"key": "value"},
	}
}
