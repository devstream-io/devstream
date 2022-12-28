package md5

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalcMD5 func", func() {
	var (
		reader io.Reader
		md5    string
		err    error

		emptyReader = strings.NewReader("")
		emtpyMd5    = "d41d8cd98f00b204e9800998ecf8427e"
		testReader  = strings.NewReader("test")
		testMd5     = "098f6bcd4621d373cade4e832627b4f6"
		errReader   io.Reader
	)

	JustBeforeEach(func() {
		md5, err = CalcMD5(reader)
	})

	When("reader is empty", func() {
		BeforeEach(func() {
			reader = emptyReader
		})
		It("should return empty md5", func() {
			Expect(md5).To(Equal(emtpyMd5))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("reader is not empty", func() {
		BeforeEach(func() {
			reader = testReader
		})
		It("should return md5", func() {
			Expect(md5).To(Equal(testMd5))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("reader is invalid", func() {
		BeforeEach(func() {
			reader = errReader
		})
		It("should return error", func() {
			Expect(err).To(HaveOccurred())
			Expect(md5).To(BeEmpty())
		})
	})
})
