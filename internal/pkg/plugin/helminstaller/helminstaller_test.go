package helminstaller_test

import (
	"os"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("helm installer test", func() {
	Context("GetDefaultOptionsByInstanceID", func() {
		defaults.DefaultOptionsMap["foo"] = &helm.Options{
			ValuesYaml: "foo: bar",
		}

		It("should exists", func() {
			opts := helminstaller.GetDefaultOptionsByInstanceID("foo-001")
			Expect(opts).NotTo(BeNil())
			Expect(opts.ValuesYaml).To(Equal("foo: bar"))
		})

		It("should not exists", func() {
			optsNil := helminstaller.GetDefaultOptionsByInstanceID("fo-001")
			Expect(optsNil).To(BeNil())
		})
	})

	Context("GetStatusGetterFuncByInstanceID", func() {
		defaults.StatusGetterFuncMap = map[string]plugininstaller.StatusGetterOperation{
			"foo": func(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
				return nil, nil
			},
		}

		It("should exists", func() {
			fn := helminstaller.GetStatusGetterFuncByInstanceID("foo-001")
			Expect(fn).NotTo(BeNil())
		})

		It("should not exists", func() {
			fn := helminstaller.GetStatusGetterFuncByInstanceID("fooo")
			Expect(fn).To(BeNil())
		})
	})

	Context("RenderDefaultConfig", func() {
		opts := configmanager.RawOptions{}
		opts["instanceID"] = interface{}("argocd-001")
		newOpts, err := helminstaller.RenderDefaultConfig(opts)
		Expect(err).To(BeNil())

		helmOpts, err := helm.NewOptions(newOpts)
		Expect(err).To(BeNil())

		Expect(helmOpts.Chart.ChartName).To(Equal(defaults.DefaultConfigWithArgoCD.Chart.ChartName))
		Expect(helmOpts.Repo.URL).To(Equal(defaults.DefaultConfigWithArgoCD.Repo.URL))
	})

	Context("RenderValuesYaml", func() {
		It("config with yaml", func() {
			opts := configmanager.RawOptions{}
			opts["valuesYaml"] = interface{}("foo: bar")
			newOpts, err := helminstaller.RenderValuesYaml(opts)
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
			newOpts, err := helminstaller.RenderValuesYaml(opts)
			Expect(err).To(BeNil())

			helmOpts, err := helm.NewOptions(newOpts)
			Expect(err).To(BeNil())

			Expect(helmOpts.Chart.ValuesYaml).To(Equal("foo: bar"))

			err = os.RemoveAll("./values.yaml")
			Expect(err).To(BeNil())
		})
	})

	Context("IndexStatusGetterFunc", func() {
		opts1 := configmanager.RawOptions{
			"instanceID": interface{}("argocd-001"),
		}

		defaults.StatusGetterFuncMap["argocd"] = defaults.GetArgoCDStatus

		fn1 := helminstaller.IndexStatusGetterFunc(opts1)
		Expect(reflect.ValueOf(fn1).Pointer()).To(Equal(reflect.ValueOf(defaults.GetArgoCDStatus).Pointer()))
	})
})
