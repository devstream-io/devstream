package git_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("RepoInfo struct", func() {
	var (
		repoName, branch, owner, org string
		repoInfo                     *git.RepoInfo
	)
	BeforeEach(func() {
		repoName = "test_repo"
		branch = "test_branch"
		owner = "test_owner"
		org = "test_org"
		repoInfo = &git.RepoInfo{
			Repo:   repoName,
			Branch: branch,
			Owner:  owner,
			Org:    org,
		}
	})
	Context("GetRepoOwner method", func() {
		It("should return owner", func() {
			result := repoInfo.GetRepoOwner()
			Expect(result).Should(Equal(repoInfo.Org))
		})
	})
	Context("GetRepoPath method", func() {
		It("should return repo path", func() {
			result := repoInfo.GetRepoPath()
			Expect(result).Should(Equal(fmt.Sprintf("%s/%s", repoInfo.GetRepoOwner(), repoName)))
		})
	})

	Context("BuildRepoRenderConfig method", func() {
		It("should return map", func() {
			config := repoInfo.BuildRepoRenderConfig()
			appName, ok := config["AppName"]
			Expect(ok).Should(BeTrue())
			Expect(appName).Should(Equal(repoName))
			repoInfo, ok := config["Repo"]
			Expect(ok).Should(BeTrue())
			repoInfoMap := repoInfo.(map[string]string)
			repoName, ok := repoInfoMap["Name"]
			Expect(ok).Should(BeTrue())
			Expect(repoName).Should(Equal(repoName))
			owner, ok := repoInfoMap["Owner"]
			Expect(ok).Should(BeTrue())
			Expect(owner).Should(Equal(org))
		})
	})

	Context("GetRepoNameWithBranch method", func() {
		BeforeEach(func() {
			repoInfo.Repo = "test"
			repoInfo.Branch = "test_branch"
		})
		It("should return repo with branch", func() {
			Expect(repoInfo.GetRepoNameWithBranch()).Should(Equal("test-test_branch"))
		})
	})

	Context("GetBranchWithDefault method", func() {
		When("branch is not empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = "test"
			})
			It("should get branch name", func() {
				Expect(repoInfo.GetBranchWithDefault()).Should(Equal("test"))
			})
		})
		When("repo is gitlab and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "gitlab"
			})
			It("should return master branch", func() {
				branch := repoInfo.GetBranchWithDefault()
				Expect(branch).Should(Equal("master"))
			})
		})
		When("repo is github and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "github"
			})
			It("should return main branch", func() {
				branch := repoInfo.GetBranchWithDefault()
				Expect(branch).Should(Equal("main"))
			})
		})
	})

	Context("BuildScmURL method", func() {
		When("repo is empty", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "not_exist"
			})
			It("should return empty url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(BeEmpty())
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
			})
			It("should return github url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("https://github.com/%s/%s", repoInfo.Org, repoInfo.Repo)))
			})
		})
		When("repo is gitlab", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
				repoInfo.BaseURL = "http://test.com"
				repoInfo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("%s/%s/%s.git", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)))
			})
		})
		When("repo is gitlab and BaseURL is not configured", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
				repoInfo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("https://gitlab.com/%s/%s.git", repoInfo.Owner, repoInfo.Repo)))
			})
		})
	})

	Context("UpdateRepoPathByCloneURL method", func() {
		var (
			cloneURL string
		)
		When("cloneURL is http format", func() {
			When("url is valid", func() {
				BeforeEach(func() {
					cloneURL = "http://test.com/test_user/test_repo.git"
				})
				It("should update owner and repo", func() {
					err := repoInfo.UpdateRepoPathByCloneURL(cloneURL)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(repoInfo.Owner).Should(Equal("test_user"))
					Expect(repoInfo.Repo).Should(Equal("test_repo"))
				})
			})
			When("url is path is not valid", func() {
				BeforeEach(func() {
					cloneURL = "http://test.com/test_user"
				})
				It("should update owner and repo", func() {
					err := repoInfo.UpdateRepoPathByCloneURL(cloneURL)
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("git repo path is not valid"))
				})
			})
		})
		When("cloneURL is git ssh format", func() {
			When("ssh format is valid", func() {
				BeforeEach(func() {
					cloneURL = "git@test.com:devstream-io/devstream.git"
				})
				It("should update owner and repo", func() {
					err := repoInfo.UpdateRepoPathByCloneURL(cloneURL)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(repoInfo.Owner).Should(Equal("devstream-io"))
					Expect(repoInfo.Repo).Should(Equal("devstream"))
				})
			})
		})
		When("ssh format has not valid path", func() {
			BeforeEach(func() {
				cloneURL = "git@test.com"
			})
			It("should return error", func() {
				err := repoInfo.UpdateRepoPathByCloneURL(cloneURL)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("scm git ssh repo not valid"))
			})
		})
		When("cloneURL is not valid", func() {
			BeforeEach(func() {
				cloneURL = "I'm just a string"
			})
			It("should return error", func() {
				err := repoInfo.UpdateRepoPathByCloneURL(cloneURL)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("git repo transport not support for now"))
			})
		})
	})
	Context("CheckValid method", func() {
		var (
			gitlabEnv, githubEnv string
		)
		BeforeEach(func() {
			gitlabEnv = os.Getenv("GITLAB_TOKEN")
			githubEnv = os.Getenv("GITHUB_TOKEN")
			os.Unsetenv("GITLAB_TOKEN")
			os.Unsetenv("GITHUB_TOKEN")
		})
		When("gitlab token is not configured", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
			})
			It("should return err", func() {
				err := repoInfo.CheckValid()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("pipeline gitlab should set env GITLAB_TOKEN"))
			})
		})
		When("github token is not configured", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
			})
			It("should return err", func() {
				err := repoInfo.CheckValid()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("pipeline github should set env GITHUB_TOKEN"))
			})
		})
		When("token is configured", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
				os.Setenv("GITHUB_TOKEN", "test")
			})
			It("should return err", func() {
				err := repoInfo.CheckValid()
				Expect(err).Error().ShouldNot(HaveOccurred())
			})

		})

		AfterEach(func() {
			if githubEnv != "" {
				os.Setenv("GITHUB_TOKEN", githubEnv)
			}
			if gitlabEnv != "" {
				os.Setenv("GITLAB_TOKEN", gitlabEnv)
			}
		})
	})

	Context("BuildWebhookInfo method", func() {
		var (
			baseURL, appName, token string
		)
		BeforeEach(func() {
			baseURL = "test.com"
			appName = "test"
			token = "test_token"
		})
		When("repo type is gitlab", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
			})
			It("should return gitlab webhook address", func() {
				w := repoInfo.BuildWebhookInfo(baseURL, appName, token)
				Expect(w.Address).Should(Equal(fmt.Sprintf("%s/project/%s", baseURL, appName)))
				Expect(w.SecretToken).Should(Equal(token))
			})
		})
		When("repo type is github", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
			})
			It("should return github webhook address", func() {
				w := repoInfo.BuildWebhookInfo(baseURL, appName, token)
				Expect(w.Address).Should(Equal(fmt.Sprintf("%s/github-webhook/", baseURL)))
				Expect(w.SecretToken).Should(Equal(token))
			})
		})
	})

	Context("Encode method", func() {
		It("should return map", func() {
			m := repoInfo.Encode()
			Expect(m).Should(Equal(map[string]any{
				"repo":   "test_repo",
				"branch": "test_branch",
				"owner":  "test_owner",
				"org":    "test_org",
			}))
		})
	})
})
