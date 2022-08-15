package github_test

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("Workflow", func() {
	var s *ghttp.Server
	var rightClient, wrongClient *github.Client
	var wantS []string
	var wantM map[string]error
	owner, repoName, f, org := "o", "r", ".github/workflows", "or"
	branch := "b"
	wsFiles := []string{"file1", "file2"}
	rightOpt := &git.RepoInfo{
		Owner: owner,
		Repo:  repoName,
		Org:   org,
	}
	wrongOpt := &git.RepoInfo{
		Owner: owner,
		Repo:  "",
		Org:   org,
	}
	partialFilesInRemoteDir := `[{
		"type": "file",
		"encoding": "base64",
		"size": 20678,
		"name": "file1",
		"path": "LICENSE"
	  }]`
	fullFilesInRemoteDir := `[
		{
			"type": "file",
			"encoding": "base64",
			"size": 20678,
			"name": "file1",
			"path": "LICENSE"
		},
		{
			"type": "file",
			"encoding": "base64",
			"size": 20678,
			"name": "file2",
			"path": "LICENSE"
		}
	]`
	workflows := []*github.Workflow{
		{WorkflowFileName: "file1"},
		{WorkflowFileName: "file2"},
	}
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", org, repoName, f)
	u2 := fmt.Sprintf("/repos/%s/%s/contents/%s", org, repoName, ".github/workflows/"+workflows[0].WorkflowFileName)
	allFileFoundMap := map[string]error{
		"file1": nil,
		"file2": nil,
	}
	allFileNotFoundMap := map[string]error{
		"file1": fmt.Errorf("not found"),
		"file2": fmt.Errorf("not found"),
	}
	partialFileNotFoundMap := map[string]error{
		"file1": nil,
		"file2": fmt.Errorf("not found"),
	}
	BeforeEach(func() {
		s = ghttp.NewServer()
		rightClient, _ = github.NewClientWithOption(rightOpt, s.URL())
		Expect(rightClient).NotTo(Equal(nil))
		wrongClient, _ = github.NewClientWithOption(wrongOpt, s.URL())
		Expect(wrongClient).NotTo(Equal(nil))
	})

	AfterEach(func() {
		// shut down the server between tests
		s.Close()
	})

	Describe("AddWorkflow", func() {
		It("got sha and WorkflowFileName already exists", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": "1212"}`)
			})
			err := rightClient.AddWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
		})
		It("got getFileSHA error", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": "1212"}`)
			})
			err := wrongClient.AddWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring("500"))
		})
		It("WorkflowFileName not exist and add one failed", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			err := wrongClient.AddWorkflow(workflows[0], branch)
			Expect(err).NotTo(Succeed())
		})
		It("WorkflowFileName not exist and add one successfully", func() {
			s.Reset()
			s.RouteToHandler("PUT", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": ""}`)
			})
			err := rightClient.AddWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(2))
			Expect(err).To(Succeed())
		})
	})

	Describe("DeleteWorkflow", func() {
		It("got sha and WorkflowFileName ok", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": ""}`)
			})
			err := rightClient.DeleteWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
		})
		It("got getFileSHA error", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			err := wrongClient.DeleteWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring("500"))
		})
		It("WorkflowFileName exist but delete it failed", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": "213"}`)
			})
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			err := rightClient.DeleteWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(2))
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring("404"))
		})
		It("WorkflowFileName exists and delete it successfully", func() {
			s.Reset()
			s.RouteToHandler("DELETE", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			s.RouteToHandler("GET", github.BaseURLPath+u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"sha": "213"}`)
			})
			err := rightClient.DeleteWorkflow(workflows[0], branch)
			Expect(s.ReceivedRequests()).Should(HaveLen(2))
			Expect(err).To(Succeed())
		})
	})

	Describe("FetchRemoteContent", func() {
		It("got GetContents error is both not equal nil and 404", func() {
			wsFiles := []string{"file1", "file2"}
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotImplemented)
			s, m, err := wrongClient.FetchRemoteContent(wsFiles)
			Expect(s).To(Equal(wantS))
			Expect(m).To(Equal(wantM))
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusNotImplemented)))
		})
		It("got GetContents error is not equal nil and equal 404", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			s, m, err := wrongClient.FetchRemoteContent(wsFiles)
			var wantS []string
			Expect(s).To(Equal(wantS))
			Expect(m).To(Equal(allFileNotFoundMap))
			Expect(err).To(Succeed())
		})
		It("got GetContents error is nil and status is not 200", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{}`)
			})
			s, m, err := rightClient.FetchRemoteContent(wsFiles)
			Expect(s).To(Equal(wantS))
			Expect(m).To(Equal(wantM))
			Expect(err.Error()).To(ContainSubstring("500"))
		})
		It("got GetContents error is nil and status is 200", func() {
			wantS := []string{"file1", "file2"}
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[
					{
						"type": "dir",
						"name": "file1",
						"path": "lib"
					},
				  	{
						"type": "file",
						"size": 20678,
						"name": "file2",
						"path": "LICENSE"
				  	}
				]`)
			})
			s, m, err := rightClient.FetchRemoteContent(wsFiles)
			Expect(s).To(Equal(wantS))
			Expect(m).To(Equal(wantM))
			Expect(err).To(Succeed())
		})

	})

	Describe("fetching VerifyWorkflows", func() {
		It("not found and return errMap", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			res, err := wrongClient.VerifyWorkflows(workflows)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(allFileNotFoundMap))
		})
		It("error is not equal nil and status is not 404", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotImplemented)
			res, err := wrongClient.VerifyWorkflows(workflows)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusNotImplemented)))
			var m map[string]error
			Expect(res).To(Equal(m))
		})
		It("status is 200 and some files lost", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprint(w, partialFilesInRemoteDir)
			})
			res, err := rightClient.VerifyWorkflows(workflows)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(partialFileNotFoundMap))
		})
		It("status is 200 and no files lost", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprint(w, fullFilesInRemoteDir)
			})
			res, err := rightClient.VerifyWorkflows(workflows)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(allFileFoundMap))
		})
	})

	Describe("fetching GetWorkflowPath", func() {
		It("not found", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			res, err := wrongClient.GetWorkflowPath()
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(""))
		})
		It("status is both not equal 200 and 404", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotImplemented)
			res, err := wrongClient.GetWorkflowPath()
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).NotTo(Succeed())
			Expect(res).To(Equal(""))
		})
		It("200", func() {
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprint(w, "{}")
			})

			res, err := rightClient.GetWorkflowPath()
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(github.BaseURLPath + u))
		})
	})

})
