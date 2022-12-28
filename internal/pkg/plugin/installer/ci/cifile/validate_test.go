package cifile_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("validate func", func() {
	var option configmanager.RawOptions
	When("options is not valid", func() {
		It("should return error", func() {
			ciTypeNotExistOptions := map[string]any{
				"ci": map[string]any{
					"configLocation": "workflows/Jenkinsfile",
					"type":           "gg",
				},
				"projectRepo": map[string]any{
					"baseURL": "http://127.0.0.1:30020",
					"branch":  "main",
					"org":     "",
					"owner":   "test_user",
					"name":    "test",
					"scmType": "gitlab",
				},
			}
			_, err := cifile.Validate(ciTypeNotExistOptions)
			Expect(err).Should(HaveOccurred())
		})
	})
	When("options is valid", func() {
		BeforeEach(func() {
			option = map[string]any{
				"ci": map[string]any{
					"configLocation": "workflows/Jenkinsfile",
					"type":           "jenkins",
				},
				"projectRepo": map[string]any{
					"baseURL": "http://127.0.0.1:30020",
					"branch":  "main",
					"org":     "",
					"owner":   "test_user",
					"name":    "test",
					"scmType": "gitlab",
				},
			}
		})
		It("should return nil error", func() {
			option = map[string]any{
				"ci": map[string]any{
					"configLocation": "http://test.com",
					"type":           "gitlab",
				},
				"scm": map[string]any{
					"baseURL": "http://127.0.0.1:30020",
					"branch":  "main",
					"org":     "",
					"owner":   "test_user",
					"name":    "test",
					"scmType": "gitlab",
				},
			}
			_, err := cifile.Validate(option)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("SetDefaultConfig func", func() {
	var defaultOpts *cifile.Options
	BeforeEach(func() {
		defaultCIFileConfig := &cifile.CIFileConfig{
			Type:           "github",
			ConfigLocation: "http://www.test.com",
		}
		defaultRepo := &git.RepoInfo{
			Owner:    "test",
			Repo:     "test_repo",
			Branch:   "test_branch",
			RepoType: "gitlab",
		}
		defaultOpts = &cifile.Options{
			CIFileConfig: defaultCIFileConfig,
			ProjectRepo:  defaultRepo,
		}
	})
	It("should work normal", func() {
		defaultFunc := cifile.SetDefaultConfig(defaultOpts)
		rawOptions := map[string]interface{}{}
		opts, err := defaultFunc(rawOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(opts)).Should(Equal(2))
	})
})
