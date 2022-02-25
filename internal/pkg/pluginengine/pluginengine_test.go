package pluginengine_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginengine"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var _ = Describe("Pluginengine", func() {
	var (
		smgr statemanager.Manager
		err  error
	)

	BeforeEach(func() {
		defer GinkgoRecover()

		smgr, err = statemanager.NewManager()
		Expect(err).ToNot(HaveOccurred())
		Expect(smgr).NotTo(BeNil())
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
		changes, _ := pluginengine.GetChangesForApply(smgr, cfg)

		Expect(len(changes)).To(Equal(1))
		c := changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Plugin.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionCreate))
	})

	It("should be 'two install'", func() {
		name1, name2 := "a", "b"
		kind1, kind2 := "tool-a", "too-b"
		version1, version2 := "v0.0.1", "v0.0.2"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name1, kind1, version1), *getTool(name2, kind2, version2)},
		}
		changes, _ := pluginengine.GetChangesForApply(smgr, cfg)

		Expect(len(changes)).To(Equal(2))
		c1 := changes[0]
		Expect(c1.Tool.Name).To(Equal(name1))
		Expect(c1.Tool.Plugin.Kind).To(Equal(kind1))
		Expect(c1.Tool.Plugin.Version).To(Equal(version1))
		Expect(c1.ActionName).To(Equal(statemanager.ActionCreate))

		c2 := changes[1]
		Expect(c2.Tool.Name).To(Equal(name2))
		Expect(c2.Tool.Plugin.Kind).To(Equal(kind2))
		Expect(c2.Tool.Plugin.Version).To(Equal(version2))
		Expect(c2.ActionName).To(Equal(statemanager.ActionCreate))
	})

	It("should be 1 uninstall when `dtm delete` is triggered against a config with 1 tool and a successful state", func() {
		name := "a"
		kind := "tool-a"
		version := "v0.0.1"

		cfg := &configloader.Config{
			Tools: []configloader.Tool{*getTool(name, kind, version)},
		}

		err = smgr.AddState(fmt.Sprintf("%s_%s", name, kind), statemanager.State{})
		Expect(err).NotTo(HaveOccurred())
		changes, _ := pluginengine.GetChangesForDelete(smgr, cfg)

		Expect(len(changes)).To(Equal(1))
		c := changes[0]
		Expect(c.Tool.Name).To(Equal(name))
		Expect(c.Tool.Plugin.Kind).To(Equal(kind))
		Expect(c.Tool.Plugin.Version).To(Equal(version))
		Expect(c.ActionName).To(Equal(statemanager.ActionDelete))
	})
})

func getTool(name, kind, version string) *configloader.Tool {
	return &configloader.Tool{
		Name:    name,
		Plugin:  configloader.Plugin{Kind: kind, Version: version},
		Options: map[string]interface{}{"key": "value"},
	}
}
