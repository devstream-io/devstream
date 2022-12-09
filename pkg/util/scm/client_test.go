package scm_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("PushInitRepo func", func() {
	var (
		mockRepo   *scm.MockScmClient
		commitInfo *git.CommitInfo
		err        error
	)
	BeforeEach(func() {
		commitInfo = &git.CommitInfo{
			CommitMsg:    "test",
			CommitBranch: "test-branch",
		}
	})
	When("init method return err", func() {
		BeforeEach(func() {
			mockRepo = &scm.MockScmClient{
				InitRaiseError: errors.New("init error"),
			}
		})
		It("should return err", func() {
			err = scm.PushInitRepo(mockRepo, commitInfo)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("init error"))
		})
	})

	When("push method failed", func() {
		BeforeEach(func() {
			mockRepo = &scm.MockScmClient{
				PushRaiseError: errors.New("push error"),
				NeedRollBack:   false,
			}
		})
		It("should return err", func() {
			err = scm.PushInitRepo(mockRepo, commitInfo)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("push error"))
			Expect(mockRepo.DeleteFuncIsRun).Should(BeFalse())
		})

		When("push method return needRollBack", func() {
			BeforeEach(func() {
				mockRepo = &scm.MockScmClient{
					PushRaiseError: errors.New("push error"),
					NeedRollBack:   true,
				}
			})
			It("should run DeleteRepo method", func() {
				err = scm.PushInitRepo(mockRepo, commitInfo)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("push error"))
				Expect(mockRepo.DeleteFuncIsRun).Should(BeTrue())
			})
		})
	})
})

var _ = Describe("NewClientWithAuth func", func() {
	var (
		r *git.RepoInfo
	)
	BeforeEach(func() {
		r = &git.RepoInfo{}
	})
	When("scm type not valid", func() {
		BeforeEach(func() {
			r.RepoType = "not_exist"
		})
		It("should return error", func() {
			_, err := scm.NewClientWithAuth(r)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("git scmType only support gitlab and github"))
		})
	})
})
