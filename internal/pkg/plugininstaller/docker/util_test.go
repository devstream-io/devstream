package docker

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoveDirs func", func() {

	var (
		errs []error
		dirs []string
	)

	AfterEach(func() {
		// check directories are removed
		for _, dir := range dirs {
			_, err := os.Stat(dir)
			Expect(os.IsNotExist(err)).To(BeTrue())
		}
	})

	When("the directories are not exist", func() {

		BeforeEach(func() {
			dirs = []string{"dir/not/exist/1", "dir/not/exist/2"}
		})

		It("should return no error", func() {
			// Remove the directories
			errs = RemoveDirs(dirs)
			for _, e := range errs {
				Expect(e).ToNot(HaveOccurred())
			}
		})
	})

	When("the directories are exist", func() {

		BeforeEach(func() {
			// create temp dir, it will be removed automatically after the test
			parentDir := GinkgoT().TempDir()
			dirs = []string{
				parentDir + "dir1",
				parentDir + "dir2",
				parentDir + "dir3/dir3-1",
			}
			// create directories
			for _, dir := range dirs {
				err := os.MkdirAll(dir, 0755)
				Expect(err).ToNot(HaveOccurred())
			}
			// create files
			_, err := os.CreateTemp(dirs[0], "file1*")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should remove all directories successfully", func() {
			// Remove the directories
			errs = RemoveDirs(dirs)
			for _, e := range errs {
				Expect(e).ToNot(HaveOccurred())
			}
		})
	})
})
