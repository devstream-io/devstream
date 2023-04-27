package commit_test

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/devstream-io/devstream/internal/pkg/commit"
)

var _ = Describe("Commit", func() {
	var testRepoDir string

	BeforeEach(func() {
		// 1. Create a temporary directory
		var err error
		testRepoDir, err = os.MkdirTemp("", "test-repo-*")
		Expect(err).NotTo(HaveOccurred())

		// 2. Change the working directory to the temporary directory
		err = os.Chdir(testRepoDir)
		Expect(err).NotTo(HaveOccurred())

		// 3. Initialize a git repository
		cmd := exec.Command("git", "init")
		err = cmd.Run()
		Expect(err).NotTo(HaveOccurred())

		// 4. Create a file and write some content to it
		file, err := os.Create(filepath.Join(testRepoDir, "test.txt"))
		Expect(err).NotTo(HaveOccurred())

		_, err = file.WriteString("Test content")
		Expect(err).NotTo(HaveOccurred())
		file.Close()

		// 5. Add the file to the git index
		cmd = exec.Command("git", "add", "test.txt")
		err = cmd.Run()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(testRepoDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should create a new commit with the given message", func() {
		message := "Test commit"
		err := Commit(message)
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command("git", "log", "--oneline")
		output, err := cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		Expect(string(output)).To(ContainSubstring(message))
	})

	It("should return an error when git is not installed", func() {
		origGitPath, err := exec.LookPath("git")
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv("PATH", "")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			err = os.Setenv("PATH", origGitPath)
			Expect(err).NotTo(HaveOccurred())
		}()

		err = Commit("Test commit")
		Expect(err).To(HaveOccurred())
	})
})
