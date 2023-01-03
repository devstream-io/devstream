package helminstaller

import (
	"os"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
)

var _ = Describe("helm installer test", func() {

	Context("renderDefaultConfig", func() {
		opts := configmanager.RawOptions{}
		opts["instanceID"] = interface{}("argocd-001")
		newOpts, err := renderDefaultConfig(opts)
		Expect(err).To(BeNil())

		helmOpts, err := helm.NewOptions(newOpts)
		Expect(err).To(BeNil())

		Expect(helmOpts.Chart.ChartName).To(Equal(defaults.DefaultConfigWithArgoCD.Chart.ChartName))
		Expect(helmOpts.Repo.URL).To(Equal(defaults.DefaultConfigWithArgoCD.Repo.URL))
	})

	Context("renderValuesYaml", func() {
		It("config with yaml", func() {
			opts := configmanager.RawOptions{}
			opts["valuesYaml"] = interface{}("foo: bar")
			newOpts, err := renderValuesYaml(opts)
			Expect(err).To(BeNil())

			helmOpts, err := helm.NewOptions(newOpts)
			Expect(err).To(BeNil())

			Expect(helmOpts.Chart.ValuesYaml).To(Equal("foo: bar"))
		})

		It("config with file path", func() {
			err := os.WriteFile("./values.yaml", []byte("foo: bar"), 0644)
			Expect(err).To(BeNil())

			opts := configmanager.RawOptions{}
			opts["valuesYaml"] = interface{}("./values.yaml")
			newOpts, err := renderValuesYaml(opts)
			Expect(err).To(BeNil())

			helmOpts, err := helm.NewOptions(newOpts)
			Expect(err).To(BeNil())

			Expect(helmOpts.Chart.ValuesYaml).To(Equal("foo: bar"))

			err = os.RemoveAll("./values.yaml")
			Expect(err).To(BeNil())
		})
	})

	Context("indexStatusGetterFunc", func() {
		opts1 := configmanager.RawOptions{
			"instanceID": interface{}("argocd-001"),
		}

		fn1 := indexStatusGetterFunc(opts1)
		Expect(reflect.ValueOf(fn1).Pointer()).To(Equal(reflect.ValueOf(defaults.GetArgoCDStatus).Pointer()))
	})
})
