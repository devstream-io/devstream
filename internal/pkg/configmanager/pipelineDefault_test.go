package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("jenkinsGenerator func", func() {
	var (
		jenkinsOptions  RawOptions
		pipelineOptions *pipelineGlobalOption
	)
	BeforeEach(func() {
		jenkinsOptions = RawOptions{
			"jenkins": RawOptions{
				"url": "test.jenkins.com",
			},
			"imageRepo": RawOptions{
				"user": "test_user",
			},
		}
		pipelineOptions = &pipelineGlobalOption{
			AppSpec: &appSpec{
				Language:  "test_language",
				FrameWork: "test_framework",
			},
			Repo: &git.RepoInfo{
				CloneURL: "test.scm.com",
			},
		}
	})
	It("should config default options", func() {
		opt := jenkinsGenerator(jenkinsOptions, pipelineOptions)
		Expect(opt).Should(Equal(RawOptions{
			"pipeline": RawOptions{
				"language": RawOptions{
					"name":      "test_language",
					"framework": "test_framework",
				},
				"jenkins": RawOptions{
					"url": "test.jenkins.com",
				},
				"imageRepo": RawOptions{
					"user": "test_user",
				},
			},
			"scm": RawOptions{
				"url": git.ScmURL("test.scm.com"),
			},
			"jenkins": RawOptions{
				"url": "test.jenkins.com",
			},
		}))
	})
})

var _ = Describe("pipelineArgocdAppGenerator func", func() {
	var (
		options      RawOptions
		pipelineVars *pipelineGlobalOption
	)
	BeforeEach(func() {
		options = RawOptions{}
		pipelineVars = &pipelineGlobalOption{
			ImageRepo: map[string]any{
				"owner": "test_user",
			},
			Repo: &git.RepoInfo{
				CloneURL: "scm.test.com",
			},
			AppName: "test_app",
		}
	})
	It("should return options", func() {
		opt := pipelineArgocdAppGenerator(options, pipelineVars)
		Expect(opt).Should(Equal(RawOptions{
			"app": RawOptions{
				"name":      "test_app",
				"namespace": "argocd",
			},
			"destination": RawOptions{
				"server":    "https://kubernetes.default.svc",
				"namespace": "default",
			},
			"source": RawOptions{
				"valuefile": "values.yaml",
				"path":      "helm/test_app",
				"repoURL":   "scm.test.com",
			},
			"imageRepo": RawOptions{
				"owner": "test_user",
			},
		}))
	})
})
