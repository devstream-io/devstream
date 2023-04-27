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
		// 1. Create a temporary directory
		var err error
		tempDir, err = os.MkdirTemp("", "patcher-tests")
		Expect(err).NotTo(HaveOccurred())

		// 2. Change the working directory to the temporary directory
		err = os.Chdir(tempDir)
		Expect(err).NotTo(HaveOccurred())

		// 3. Create the original file and the patch file
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

			err = Patch(patchFile.Name())
			Expect(err).NotTo(HaveOccurred())

			patchedContent, err := os.ReadFile(originalFile.Name())
			Expect(err).NotTo(HaveOccurred())

			expectedPatchedContent := `Hello, world!
This is the patched file.
`
			patchedContentStr := string(patchedContent)
			Expect(patchedContentStr).To(Equal(expectedPatchedContent))
		})

		It("returns an error if the patch file is invalid", func() {
			originalContent := `Hello, world!
This is the original file.
`

			err := os.WriteFile(originalFile.Name(), []byte(originalContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			invalidPatchContent := fmt.Sprintf(`--- %s
+++ new-file
@@ -1,2 +1,2 @@
`,
				filepath.Base(originalFile.Name()))

			err = os.WriteFile(patchFile.Name(), []byte(invalidPatchContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			err = Patch(patchFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "patch command failed")).To(BeTrue())
		})
	})

	Context("when patching a file with inconsistent indentation", func() {
		It("successfully applies the patch with spaces to the original file with tabs", func() {
			originalContent := "Hello, world!\n\tThis is the original file with tabs.\n"

			err := os.WriteFile(originalFile.Name(), []byte(originalContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			patchContent := fmt.Sprintf(`--- %s
+++ new-file
@@ -1,2 +1,2 @@
 Hello, world!
-    This is the original file with tabs.
+    This is the patched file with tabs.
`,
				filepath.Base(originalFile.Name()))

			err = os.WriteFile(patchFile.Name(), []byte(patchContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			err = Patch(patchFile.Name())
			Expect(err).NotTo(HaveOccurred())

			patchedContent, err := os.ReadFile(originalFile.Name())
			Expect(err).NotTo(HaveOccurred())

			expectedPatchedContent := "Hello, world!\n\tThis is the patched file with tabs.\n"
			Expect(string(patchedContent)).To(Equal(expectedPatchedContent))
		})

		It("successfully applies the patch with tabs to the original file with spaces", func() {
			originalContent := "Hello, world!\n    This is the original file with spaces.\n"

			err := os.WriteFile(originalFile.Name(), []byte(originalContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			patchContent := fmt.Sprintf(`--- %s
+++ new-file
@@ -1,2 +1,2 @@
 Hello, world!
-	This is the original file with spaces.
+	This is the patched file with spaces.
`,
				filepath.Base(originalFile.Name()))

			err = os.WriteFile(patchFile.Name(), []byte(patchContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			err = Patch(patchFile.Name())
			Expect(err).NotTo(HaveOccurred())

			patchedContent, err := os.ReadFile(originalFile.Name())
			Expect(err).NotTo(HaveOccurred())

			expectedPatchedContent := "Hello, world!\n    This is the patched file with spaces.\n"
			Expect(string(patchedContent)).To(Equal(expectedPatchedContent))
		})
	})
})
