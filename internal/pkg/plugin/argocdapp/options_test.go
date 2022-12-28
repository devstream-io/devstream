package argocdapp

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("source struct", func() {
	var (
		scmClient scm.ClientOperation
		errMsg    string
		s         *source
	)
	BeforeEach(func() {
		s = &source{
			Path: "test_path",
		}
	})
	Context("checkPathExist method", func() {
		When("scm client error", func() {
			BeforeEach(func() {
				errMsg = "scm get path error"
				scmClient = &scm.MockScmClient{
					GetPathInfoError: errors.New(errMsg),
				}
			})
			It("should return error", func() {
				_, err := s.checkPathExist(scmClient)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("path list is empty", func() {
			BeforeEach(func() {
				scmClient = &scm.MockScmClient{
					GetPathInfoReturnValue: []*git.RepoFileStatus{},
				}
			})
			It("should return false", func() {
				exist, err := s.checkPathExist(scmClient)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(exist).Should(BeFalse())
			})
		})
		When("path list is not empty", func() {
			BeforeEach(func() {
				scmClient = &scm.MockScmClient{
					GetPathInfoReturnValue: []*git.RepoFileStatus{
						{Path: "gg"},
					},
				}
			})
			It("should return true", func() {
				exist, err := s.checkPathExist(scmClient)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(exist).Should(BeTrue())
			})
		})
	})
})

var _ = Describe("options struct", func() {
	var (
		configLocation    downloader.ResourceLocation
		o                 *options
		s                 *ghttp.Server
		pathName, reqPath string
	)
	BeforeEach(func() {
		s = ghttp.NewServer()
		o = &options{
			App: &app{
				Name: "test",
			},
			Source: &source{
				Path:    "test_path",
				RepoURL: "https://test.com/owner/repo",
			},
			ImageRepo: &imageRepo{
				URL:       "test.com",
				User:      "user",
				InitalTag: "latest",
			},
		}
		pathName = "/configLoc"
		reqPath = fmt.Sprintf("%s%s", s.URL(), pathName)
		s.RouteToHandler("HEAD", pathName, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
	})
	Context("getArgocdDefaultConfigFiles method", func() {
		When("can't get location", func() {
			BeforeEach(func() {
				s.RouteToHandler("GET", pathName, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", pathName),
					ghttp.RespondWith(http.StatusNotFound, ""),
				))
				configLocation = downloader.ResourceLocation(reqPath)
			})
			It("should return error", func() {
				_, err := o.getArgocdDefaultConfigFiles(configLocation)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("bad response code: 404"))
			})
		})
		When("all valid", func() {
			BeforeEach(func() {
				configFileContent := `
image:
  repository: [[ .ImageRepo.URL ]]/[[ .ImageRepo.User ]]
  tag: [[ .ImageRepo.InitalTag ]]
  pullPolicy: Always

`
				s.RouteToHandler("GET", pathName, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", pathName),
					ghttp.RespondWith(http.StatusOK, configFileContent),
				))
				configLocation = downloader.ResourceLocation(reqPath)
			})
			It("should return valid data", func() {
				data, err := o.getArgocdDefaultConfigFiles(configLocation)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).ShouldNot(BeEmpty())
			})
		})
	})
	AfterEach(func() {
		s.Close()
	})
})
