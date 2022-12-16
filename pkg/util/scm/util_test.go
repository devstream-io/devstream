package scm

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("calculateSHA func", func() {
	var content string
	It("should return as expect", func() {
		content = "test Content"
		Expect(calculateSHA([]byte(content))).Should(Equal("f73d59b513c429a33da4f7efe70c7af3"))
	})
})
