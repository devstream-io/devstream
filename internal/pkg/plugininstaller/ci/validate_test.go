package ci_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
)

var _ = Describe("Validate func", func() {
	var option map[string]interface{}
	When("options is not valid", func() {
		It("should return error", func() {
			fieldNotValidoption := map[string]interface{}{
				"name":       "ci-generic",
				"instanceID": "test",
				"options": map[string]interface{}{
					"projectRepo": map[string]interface{}{
						"owner":    "test_user",
						"org":      "",
						"repo":     "test",
						"branch":   "test",
						"repoType": "github",
					},
				},
			}

			_, err := ci.Validate(fieldNotValidoption)
			Expect(err).Should(HaveOccurred())
			fileNotExistOption := map[string]any{
				"ci": map[string]any{
					"localPath": "workflows/Jenkinsfile",
					"type":      "jenkins",
				},
				"projectRepo": map[string]any{
					"baseURL":  "http://127.0.0.1:30020",
					"branch":   "main",
					"org":      "",
					"owner":    "test_user",
					"repo":     "test",
					"repoType": "gitlab",
				},
			}
			_, err = ci.Validate(fileNotExistOption)
			Expect(err).Should(HaveOccurred())
			repoCiConflictOption := map[string]any{
				"ci": map[string]any{
					"localPath": "workflows/Jenkinsfile",
					"type":      "github",
				},
				"projectRepo": map[string]any{
					"baseURL":  "http://127.0.0.1:30020",
					"branch":   "main",
					"org":      "",
					"owner":    "test_user",
					"repo":     "test",
					"repoType": "gitlab",
				},
			}
			_, err = ci.Validate(repoCiConflictOption)
			Expect(err).Should(HaveOccurred())

			ciTypeNotExistOptions := map[string]any{
				"ci": map[string]any{
					"localPath": "workflows/Jenkinsfile",
					"type":      "gg",
				},
				"projectRepo": map[string]any{
					"baseURL":  "http://127.0.0.1:30020",
					"branch":   "main",
					"org":      "",
					"owner":    "test_user",
					"repo":     "test",
					"repoType": "gitlab",
				},
			}
			_, err = ci.Validate(ciTypeNotExistOptions)
			Expect(err).Should(HaveOccurred())

		})
	})
	When("options is valid", func() {
		BeforeEach(func() {
			option = map[string]any{
				"ci": map[string]any{
					"localPath": "workflows/Jenkinsfile",
					"type":      "jenkins",
				},
				"projectRepo": map[string]any{
					"baseURL":  "http://127.0.0.1:30020",
					"branch":   "main",
					"org":      "",
					"owner":    "test_user",
					"repo":     "test",
					"repoType": "gitlab",
				},
			}
		})
		It("should return nil error", func() {
			option = map[string]any{
				"ci": map[string]any{
					"remoteURL": "http://test.com",
					"type":      "gitlab",
				},
				"projectRepo": map[string]any{
					"baseURL":  "http://127.0.0.1:30020",
					"branch":   "main",
					"org":      "",
					"owner":    "test_user",
					"repo":     "test",
					"repoType": "gitlab",
				},
			}
			_, err := ci.Validate(option)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("SetDefaultConfig func", func() {
	var defaultOpts *ci.Options
	BeforeEach(func() {
		defaultCIConfig := &ci.CIConfig{
			Type:      "github",
			RemoteURL: "http://www.test.com",
		}
		defaultRepo := &common.Repo{
			Owner:    "test",
			Repo:     "test_repo",
			Branch:   "test_branch",
			RepoType: "gitlab",
		}
		defaultOpts = &ci.Options{
			CIConfig:    defaultCIConfig,
			ProjectRepo: defaultRepo,
		}
	})
	It("should work normal", func() {
		defaultFunc := ci.SetDefaultConfig(defaultOpts)
		rawOptions := map[string]interface{}{}
		opts, err := defaultFunc(rawOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(opts)).Should(Equal(2))
	})
})
