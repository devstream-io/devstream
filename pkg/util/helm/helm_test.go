package helm_test

import (
	"fmt"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/helm"
)

var _ = Describe("NewHelm func", func() {
	var (
		client helmclient.Client
	)
	When("params are right", func() {
		BeforeEach(func() {
			client = &mockClient{}
		})
		It("should work noraml", func() {
			got, err := helm.NewHelm(helmParam, helm.WithClient(client))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(got).ShouldNot(BeNil())
		})
	})
	When("params are wrong", func() {
		BeforeEach(func() {
			client = &mockClient{
				AddOrUpdateChartRepoError: fmt.Errorf("test error"),
			}
		})
		It("should work noraml", func() {
			got, err := helm.NewHelm(helmParam, helm.WithClient(client))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("test error"))
			Expect(got).Should(BeNil())
		})
	})
	When("has option", func() {
		var (
			entry  *repo.Entry
			atomic bool
			spec   *helmclient.ChartSpec
		)
		BeforeEach(func() {

			entry = &repo.Entry{
				Name:                  helmParam.Repo.Name,
				URL:                   helmParam.Repo.URL,
				Username:              "",
				Password:              "",
				CertFile:              "",
				KeyFile:               "",
				CAFile:                "",
				InsecureSkipTLSverify: false,
				PassCredentialsAll:    false,
			}
			atomic = true
			if !*helmParam.Chart.Wait {
				atomic = false
			}
			timeout, err := time.ParseDuration(helmParam.Chart.Timeout)
			Expect(err).ShouldNot(HaveOccurred())
			spec = &helmclient.ChartSpec{
				ReleaseName:      helmParam.Chart.ReleaseName,
				ChartName:        helmParam.Chart.ChartName,
				Namespace:        helmParam.Chart.Namespace,
				ValuesYaml:       helmParam.Chart.ValuesYaml,
				Version:          helmParam.Chart.Version,
				CreateNamespace:  false,
				DisableHooks:     false,
				Replace:          true,
				Wait:             *helmParam.Chart.Wait,
				DependencyUpdate: false,
				Timeout:          timeout,
				GenerateName:     false,
				NameTemplate:     "",
				Atomic:           atomic,
				SkipCRDs:         false,
				UpgradeCRDs:      *helmParam.Chart.UpgradeCRDs,
				SubNotes:         false,
				Force:            false,
				ResetValues:      false,
				ReuseValues:      false,
				Recreate:         false,
				MaxHistory:       0,
				CleanupOnFail:    false,
				DryRun:           false,
			}
			client = &mockClient{}
		})
		It("should work normal", func() {
			got, err := helm.NewHelm(helmParam, helm.WithClient(client))
			Expect(err).ShouldNot(HaveOccurred())
			want := &helm.Helm{
				Entry:     entry,
				ChartSpec: spec,
				Client:    client,
			}
			Expect(got).Should(Equal(want))
		})
	})
})
