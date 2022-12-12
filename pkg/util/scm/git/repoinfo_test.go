package git

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RepoInfo struct", func() {
	var (
		repoName, branch, owner, org string
		repoInfo                     *RepoInfo
	)
	BeforeEach(func() {
		repoName = "test_repo"
		branch = "test_branch"
		owner = "test_owner"
		org = "test_org"
		repoInfo = &RepoInfo{
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

	Context("GetRepoName method", func() {
		When("name field is configured", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					Repo: repoName,
				}
			})
			It("should return owner", func() {
				result := repoInfo.GetRepoName()
				Expect(result).Should(Equal(repoInfo.Repo))
			})
		})
		When("name field is not configured, url exist", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					CloneURL: ScmURL("https://test.github.com/user/urlRepo"),
				}
			})
			It("should return owner", func() {
				result := repoInfo.GetRepoName()
				Expect(result).Should(Equal("urlRepo"))
			})
		})
	})

	Context("GetCloneURL mehtod", func() {
		When("name field is configured", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					CloneURL: "exist.com",
				}
			})
			It("should return owner", func() {
				result := repoInfo.GetCloneURL()
				Expect(result).Should(Equal("exist.com"))
			})
		})
		When("url field is not configured, other fields exist", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					Repo:     "test_name",
					Owner:    "test_user",
					RepoType: "github",
				}
			})
			It("should return owner", func() {
				result := repoInfo.GetCloneURL()
				Expect(result).Should(Equal("https://github.com/test_user/test_name"))
			})
		})

	})

	Context("GetRepoPath method", func() {
		It("should return repo path", func() {
			result := repoInfo.GetRepoPath()
			Expect(result).Should(Equal(fmt.Sprintf("%s/%s", repoInfo.GetRepoOwner(), repoName)))
		})
	})

	Context("getBranchWithDefault method", func() {
		When("branch is not empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = "test"
			})
			It("should get branch name", func() {
				Expect(repoInfo.getBranchWithDefault()).Should(Equal("test"))
			})
		})
		When("repo is gitlab and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "gitlab"
			})
			It("should return master branch", func() {
				branch := repoInfo.getBranchWithDefault()
				Expect(branch).Should(Equal("master"))
			})
		})
		When("repo is github and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "github"
			})
			It("should return main branch", func() {
				branch := repoInfo.getBranchWithDefault()
				Expect(branch).Should(Equal("main"))
			})
		})
	})

	Context("buildScmURL method", func() {
		When("repo is github", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
			})
			It("should return github url", func() {
				url := repoInfo.buildScmURL()
				Expect(string(url)).Should(Equal(fmt.Sprintf("https://github.com/%s/%s", repoInfo.Org, repoInfo.Repo)))
			})
		})
		When("repo is gitlab", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
				repoInfo.BaseURL = "http://test.com"
				repoInfo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repoInfo.buildScmURL()
				Expect(string(url)).Should(Equal(fmt.Sprintf("%s/%s/%s.git", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)))
			})
		})
		When("repo is gitlab and BaseURL is not configured", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
				repoInfo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repoInfo.buildScmURL()
				Expect(string(url)).Should(Equal(fmt.Sprintf("https://gitlab.com/%s/%s.git", repoInfo.Owner, repoInfo.Repo)))
			})
		})
	})

	Context("checkValid method", func() {
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
				repoInfo = &RepoInfo{
					Owner:    "test",
					RepoType: "gitlab",
					NeedAuth: true,
					Repo:     "test",
				}
			})
			It("should return err", func() {
				err := repoInfo.checkValid()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("gitlab repo should set env GITLAB_TOKEN"))
			})
		})
		When("github token is not configured", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					Owner:    "test",
					RepoType: "github",
					NeedAuth: true,
					Repo:     "test",
				}
			})
			It("should return err", func() {
				err := repoInfo.checkValid()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("github repo should set env GITHUB_TOKEN"))
			})
		})
		When("token is configured", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					Owner:    "test",
					RepoType: "github",
					NeedAuth: true,
					Repo:     "test",
				}
				os.Setenv("GITHUB_TOKEN", "test")
			})
			It("should return err", func() {
				err := repoInfo.checkValid()
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})

		When("org and owner are all exist", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					Org:   "exist",
					Owner: "exist",
				}
			})
			It("should return error", func() {
				err := repoInfo.checkValid()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("git org and owner can't be configured at the same time"))
			})
		})

		When("repo type is not valid", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					RepoType: "not_exist",
				}
			})
			It("should return error", func() {
				err := repoInfo.checkValid()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("git scmType only support gitlab and github"))
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

	Context("Encode method", func() {
		It("should return map", func() {
			m := repoInfo.Encode()
			Expect(m).Should(Equal(map[string]any{
				"name":   "test_repo",
				"branch": "test_branch",
				"owner":  "test_owner",
				"org":    "test_org",
			}))
		})
	})

	Context("SetDefault method", func() {
		When("is github repo", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					CloneURL: "git@github.com:test/dtm-test.git",
					Branch:   "test",
				}
			})
			It("should return github repo info", func() {
				err := repoInfo.SetDefault()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repoInfo).ShouldNot(BeNil())
				Expect(repoInfo.Repo).Should(Equal("dtm-test"))
				Expect(repoInfo.Owner).Should(Equal("test"))
				Expect(repoInfo.RepoType).Should(Equal("github"))
				Expect(repoInfo.Branch).Should(Equal("test"))
			})
		})
		When("clone url is not valid", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					CloneURL: "git@github.comtest/dtm-test.git",
					Branch:   "test",
				}
			})
			It("should return error", func() {
				err := repoInfo.SetDefault()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("is gitlab repo", func() {
			When("apiURL is not set, url is ssh format", func() {
				BeforeEach(func() {
					repoInfo = &RepoInfo{
						CloneURL: "git@gitlab.test.com:root/test-demo.git",
						APIURL:   "",
						Branch:   "test",
						RepoType: "gitlab",
					}

				})
				It("should return error", func() {
					err := repoInfo.SetDefault()
					Expect(err).Error().Should(HaveOccurred())
				})
			})
			When("apiURL is not set, url is http format", func() {
				BeforeEach(func() {
					repoInfo = &RepoInfo{
						CloneURL: "http://gitlab.test.com:3000/root/test-demo.git",
						APIURL:   "",
						Branch:   "test",
						RepoType: "gitlab",
					}
				})
				It("should return error", func() {
					err := repoInfo.SetDefault()
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(repoInfo.BaseURL).Should(Equal("http://gitlab.test.com:3000"))
					Expect(repoInfo.Owner).Should(Equal("root"))
					Expect(repoInfo.Repo).Should(Equal("test-demo"))
					Expect(repoInfo.Branch).Should(Equal("test"))
				})
			})
			When("apiURL is set", func() {
				BeforeEach(func() {
					repoInfo = &RepoInfo{
						CloneURL: "git@gitlab.test.com:root/test-demo.git",
						APIURL:   "http://gitlab.http.com",
						Branch:   "test",
						RepoType: "gitlab",
						Org:      "cover_org",
					}
				})
				It("should set apiURL as BaseURL", func() {
					err := repoInfo.SetDefault()
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(repoInfo.BaseURL).Should(Equal("http://gitlab.http.com"))
					Expect(repoInfo.Owner).Should(Equal(""))
					Expect(repoInfo.Repo).Should(Equal("test-demo"))
					Expect(repoInfo.Branch).Should(Equal("test"))
					Expect(repoInfo.Org).Should(Equal("cover_org"))
				})

			})
		})
		When("scm repo has fields", func() {
			BeforeEach(func() {
				repoInfo = &RepoInfo{
					CloneURL: "https://github.com/test_org/test",
					Repo:     "test",
					RepoType: "github",
					Org:      "test_org",
				}
			})
			It("should return repoInfo", func() {
				err := repoInfo.SetDefault()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repoInfo).Should(Equal(&RepoInfo{
					Branch:   "main",
					Repo:     "test",
					RepoType: "github",
					NeedAuth: false,
					CloneURL: "https://github.com/test_org/test",
					Org:      "test_org",
				}))
			})
		})
	})
})

var _ = Describe("ScmURL type", func() {
	var cloneURL ScmURL
	Context("UpdateRepoPathByCloneURL method", func() {
		When("cloneURL is http format", func() {
			When("url is valid", func() {
				BeforeEach(func() {
					cloneURL = "http://test.com/test_user/test_repo.git"
				})
				It("should update owner and repo", func() {
					owner, name, err := cloneURL.extractRepoOwnerAndName()
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(owner).Should(Equal("test_user"))
					Expect(name).Should(Equal("test_repo"))
				})
			})
			When("url is path is not valid", func() {
				BeforeEach(func() {
					cloneURL = "http://test.com/test_user"
				})
				It("should update owner and repo", func() {
					_, _, err := cloneURL.extractRepoOwnerAndName()
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("git url repo path is not valid"))
				})
			})
		})
		When("cloneURL is git ssh format", func() {
			When("ssh format is valid", func() {
				BeforeEach(func() {
					cloneURL = "git@test.com:devstream-io/devstream.git"
				})
				It("should update owner and repo", func() {
					owner, name, err := cloneURL.extractRepoOwnerAndName()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(owner).Should(Equal("devstream-io"))
					Expect(name).Should(Equal("devstream"))
				})
			})
		})
		When("ssh format has not valid path", func() {
			BeforeEach(func() {
				cloneURL = "git@test.com"
			})
			It("should return error", func() {
				_, _, err := cloneURL.extractRepoOwnerAndName()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("git url ssh repo not valid"))
			})
		})
		When("cloneURL is not valid", func() {
			BeforeEach(func() {
				cloneURL = "I'm just a string"
			})
			It("should return error", func() {
				_, _, err := cloneURL.extractRepoOwnerAndName()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("git url repo transport not support for now"))
			})
		})
	})
	Context("addGithubURL method", func() {
		BeforeEach(func() {
			cloneURL = "github.com/test/repo"
		})
		It("should add scheme", func() {
			Expect(string(cloneURL.addGithubURLScheme())).Should(Equal("https://github.com/test/repo"))
		})
	})
})
