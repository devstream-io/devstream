package ci

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ci walkDir methods", func() {
	type testCase struct {
		filePath string
		ciType   ciRepoType
		isDir    bool
	}
	var (
		testCases    []*testCase
		testCaseData *testCase
	)
	Context("filterCIFilesFunc func", func() {
		When("input file not not right", func() {
			BeforeEach(func() {
				testCases = []*testCase{
					{"Jenkinsfile", ciRepoType("gitlab"), false},
					{".gitlab.ci", ciRepoType("jenkins"), false},
					{"Jenkinsfile", ciRepoType("jenkins"), true},
					{"workflows/pr.yaml", ciRepoType("jenkins"), false},
				}
			})
			It("should return false", func() {
				for _, tt := range testCases {
					Expect(filterCIFilesFunc(tt.ciType)(tt.filePath, tt.isDir)).Should(BeFalse())
				}
			})
		})
		When("input file is valid", func() {
			BeforeEach(func() {
				testCases = []*testCase{
					{"Jenkinsfile", ciRepoType("jenkins"), false},
					{".gitlab-ci.yml", ciRepoType("gitlab"), false},
					{"workflows/pr.yaml", ciRepoType("github"), false},
					{"workflows/pr2.yaml", ciRepoType("github"), false},
				}
			})
			It("should return true", func() {
				for _, tt := range testCases {
					Expect(filterCIFilesFunc(tt.ciType)(tt.filePath, tt.isDir)).Should(BeTrue())
				}
			})
		})
	})
	Context("processCIFilesFunc func", func() {
		var (
			tempFileLoc string
			testContent []byte
		)
		BeforeEach(func() {
			tempDir := GinkgoT().TempDir()
			tempFile, err := os.CreateTemp(tempDir, "testFile")
			Expect(err).Error().ShouldNot(HaveOccurred())
			tempFileLoc = tempFile.Name()
		})
		When("file is right", func() {
			BeforeEach(func() {
				testContent = []byte("[[ .Name ]]_test")
				err := os.WriteFile(tempFileLoc, testContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should work as expected", func() {
				result, err := processCIFilesFunc("test", map[string]interface{}{
					"Name": "devstream",
				})(tempFileLoc)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(result).Should(Equal([]byte("devstream_test")))
			})
		})
		When("file is not exist", func() {
			It("should return error", func() {
				_, err := processCIFilesFunc("test", map[string]interface{}{
					"Name": "devstream",
				})("not_exist_file")
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("getGitNameFunc func", func() {
		When("ciType is github", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"workflows/pr.yaml", ciRepoType("github"), false}
			})
			It("should return github workflows path", func() {
				result := getGitNameFunc(testCaseData.ciType)(testCaseData.filePath, "workflows")
				Expect(result).Should(Equal(".github/workflows/pr.yaml"))
			})
		})
		When("ciType is others", func() {
			BeforeEach(func() {
				testCaseData = &testCase{"work/Jenkinsfile", ciRepoType("jenkins"), false}
			})
			It("should return github workflows path", func() {
				result := getGitNameFunc(testCaseData.ciType)(testCaseData.filePath, "")
				Expect(result).Should(Equal("Jenkinsfile"))
			})
		})
	})
})
