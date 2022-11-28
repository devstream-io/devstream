package helminstaller

import (
	"os"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("helm installer test", func() {
	Context("getDefaultOptionsByInstanceID", func() {
		defaults.DefaultOptionsMap["foo"] = &helm.Options{
			ValuesYaml: "foo: bar",
		}

		It("should exists", func() {
			opts := getDefaultOptionsByInstanceID("foo-001")
			Expect(opts).NotTo(BeNil())
			Expect(opts.ValuesYaml).To(Equal("foo: bar"))
		})

		It("should not exists", func() {
			optsNil := getDefaultOptionsByInstanceID("fo-001")
			Expect(optsNil).To(BeNil())
		})
	})

	Context("getStatusGetterFuncByInstanceID", func() {
		defaults.StatusGetterFuncMap = map[string]installer.StatusGetterOperation{
			"foo": func(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
				return nil, nil
			},
		}

		It("should exists", func() {
			fn := getStatusGetterFuncByInstanceID("foo-001")
			Expect(fn).NotTo(BeNil())
		})

		It("should not exists", func() {
			fn := getStatusGetterFuncByInstanceID("fooo")
			Expect(fn).To(BeNil())
		})
	})

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

		defaults.StatusGetterFuncMap["argocd"] = defaults.GetArgoCDStatus

		fn1 := indexStatusGetterFunc(opts1)
		Expect(reflect.ValueOf(fn1).Pointer()).To(Equal(reflect.ValueOf(defaults.GetArgoCDStatus).Pointer()))
	})

	Context("getLongestMatchedName", func() {
		testList := []string{"abc", "abcd", "ab"}
		retStr := getLongestMatchedName(testList)
		Expect(retStr).To(Equal("abcd"))
	})

	Context("getDefaultOptionsByInstanceID", func() {
		opt1 := getDefaultOptionsByInstanceID("argocd")
		opt2 := getDefaultOptionsByInstanceID("argocd-001")
		opt3 := getDefaultOptionsByInstanceID("argocd001")

		Expect(opt1).To(Equal(opt2))
		Expect(opt2).To(Equal(opt3))
	})
})
