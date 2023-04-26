package patch_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/devstream-io/devstream/internal/pkg/patch"
)

func TestPatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Patcher Suite")
}

var _ = Describe("Patcher", func() {
	var (
		originalFile *os.File
		patchFile    *os.File
		tempDir      string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "patcher-tests")
		Expect(err).NotTo(HaveOccurred())

		originalFile, err = os.CreateTemp(tempDir, "original-*")
		Expect(err).NotTo(HaveOccurred())

		patchFile, err = os.CreateTemp(tempDir, "patch-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when patching a file", func() {
		It("successfully applies the patch to the original file", func() {
			originalContent := `Hello, world!
This is the original file.
`

			err := os.WriteFile(originalFile.Name(), []byte(originalContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			patchContent := fmt.Sprintf(`--- %s
+++ new-file
@@ -1,2 +1,2 @@
 Hello, world!
-This is the original file.
+This is the patched file.
`,
				filepath.Base(originalFile.Name()))

			err = os.WriteFile(patchFile.Name(), []byte(patchContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			err = Patch(tempDir, patchFile.Name())
			Expect(err).NotTo(HaveOccurred())

			patchedContent, err := os.ReadFile(originalFile.Name())
			Expect(err).NotTo(HaveOccurred())

			expectedPatchedContent := `Hello, world!
This is the patched file.
`
			Expect(string(patchedContent)).To(Equal(expectedPatchedContent))
		})

		It("returns an error if the patch command is not found or not executable", func() {
			// Temporarily change PATH to exclude the real patch command
			originalPath := os.Getenv("PATH")
			err := os.Setenv("PATH", tempDir)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				err := os.Setenv("PATH", originalPath)
				Expect(err).NotTo(HaveOccurred())
			}()

			err = Patch(tempDir, patchFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "patch command not found")).To(BeTrue())
		})

		It("returns an error if the patch file is invalid", func() {
			invalidPatchContent := `This is not a valid patch file.`
			err := os.WriteFile(patchFile.Name(), []byte(invalidPatchContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			err = Patch(tempDir, patchFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "patch command failed")).To(BeTrue())
		})
	})
})
