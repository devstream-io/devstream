package planmanager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/backend"
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
		Ω(err).ShouldNot(HaveOccurred())
		smgr = statemanager.NewManager(b)
		_, _ = GinkgoWriter.Write([]byte("new a statemanager"))
	})

	Describe("Generating plans", func() {
		Context("With config only", func() {
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
				name1 := "tool_a"
				version1 := "v0.0.1"
				name2 := "tool_b"
				version2 := "v0.0.2"
				cfg := &configloader.Config{
					Tools: []configloader.Tool{*getTool(name1, version1), *getTool(name2, version2)},
				}

				plan := planmanager.NewPlan(smgr, cfg)
				Expect(len(plan.Changes)).To(Equal(2))
				c1 := plan.Changes[0]
				c2 := plan.Changes[1]

				Expect(c1.Tool.Name).To(Equal(name1))
				Expect(c1.Tool.Version).To(Equal(version1))
				Expect(c1.ActionName).To(Equal(statemanager.ActionInstall))

				Expect(c2.Tool.Name).To(Equal(name2))
				Expect(c2.Tool.Version).To(Equal(version2))
				Expect(c2.ActionName).To(Equal(statemanager.ActionInstall))
			})
		})
	})
})

func getTool(name, version string) *configloader.Tool {
	return &configloader.Tool{
		Name:    name,
		Version: version,
		Options: map[string]interface{}{"key": "value"},
	}
}
