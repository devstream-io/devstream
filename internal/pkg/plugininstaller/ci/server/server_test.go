package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ci walkDir methods", func() {
	type testCase struct {
		filePath string
		ciType   CIServerType
		isDir    bool
	}
	var (
		testCases    []*testCase
		testCaseData *testCase
	)
	Context("FilterCIFilesFunc func", func() {
		When("input file not not right", func() {
			BeforeEach(func() {
				testCases = []*testCase{
					{"Jenkinsfile", CIServerType("gitlab"), false},
					{"directory", CIServerType("gitlab"), true},
					{"directory", CIServerType("github"), true},
					{".gitlab.ci", CIServerType("jenkins"), false},
					{"Jenkinsfile", CIServerType("jenkins"), true},
					{"workflows/pr.yaml", CIServerType("jenkins"), false},
				}
			})
			It("should return false", func() {
				for _, tt := range testCases {
					Expect(NewCIServer(tt.ciType).FilterCIFilesFunc()(tt.filePath, tt.isDir)).Should(BeFalse())
				}
			})
		})
		When("input file is valid", func() {
			BeforeEach(func() {
				testCases = []*testCase{
					{"Jenkinsfile", CIServerType("jenkins"), false},
					{".gitlab-ci.yml", CIServerType("gitlab"), false},
					{"workflows/pr.yaml", CIServerType("github"), false},
					{"workflows/pr2.yaml", CIServerType("github"), false},
				}
			})
			It("should return true", func() {
				for _, tt := range testCases {
					Expect(NewCIServer(tt.ciType).FilterCIFilesFunc()(tt.filePath, tt.isDir)).Should(BeTrue())
				}
			})
		})
	})

	Context("GetGitNameFunc func", func() {
		When("ciType is github", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"workflows/pr.yaml", CIServerType("github"), false}
			})
			It("should return github workflows path", func() {
				result := NewCIServer(testCaseData.ciType).GetGitNameFunc()(testCaseData.filePath, "workflows")
				Expect(result).Should(Equal(".github/workflows/pr.yaml"))
			})
		})
		When("ciType is jenkins", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"work/Jenkinsfile", CIServerType("jenkins"), false}
			})
			It("should return github workflows path", func() {
				result := NewCIServer(testCaseData.ciType).GetGitNameFunc()(testCaseData.filePath, "")
				Expect(result).Should(Equal("Jenkinsfile"))
			})
		})
		When("ciType is gitlab", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"work/.gitlab-ci.yml", CIServerType("gitlab"), false}
			})
			It("should return github workflows path", func() {
				result := NewCIServer(testCaseData.ciType).GetGitNameFunc()(testCaseData.filePath, "")
				Expect(result).Should(Equal(".gitlab-ci.yml"))
			})
		})
	})

	Context("CIFilePath func", func() {
		When("ciType is github", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"workflows/pr.yaml", CIServerType("github"), false}
			})
			It("should return github workflows path", func() {
				result := NewCIServer(testCaseData.ciType).CIFilePath()
				Expect(result).Should(Equal(".github/workflows/app.yaml"))
			})
		})
		When("ciType is gitlab", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"", CIServerType("gitlab"), false}
			})
			It("should return gitlab ci path", func() {
				result := NewCIServer(testCaseData.ciType).CIFilePath()
				Expect(result).Should(Equal(".gitlab-ci.yml"))
			})
		})
		When("ciType is jenkins", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"", CIServerType("jenkins"), false}
			})
			It("should return jenkins ci path", func() {
				result := NewCIServer(testCaseData.ciType).CIFilePath()
				Expect(result).Should(Equal("Jenkinsfile"))
			})
		})
	})

	Context("NewCIServer func", func() {
		When("ci type not exist", func() {
			It("should return nil", func() {
				server := NewCIServer(CIServerType("not_exist_type"))
				Expect(server).Should(BeNil())
			})
		})
	})
})
