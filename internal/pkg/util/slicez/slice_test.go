package slicez

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mapz", func() {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{1, 2, 3, 4}
	s3 := []int{1, 2, 3, 4, 5}
	s4 := []int{1, 2, 3}

	It("should be a slice with 0 items", func() {
		retSlice1 := SliceInSlice(s1, s2)
		Expect(len(retSlice1.([]int)) == 0)
		retSlice2 := SliceInSlice(s1, s3)
		Expect(len(retSlice2.([]int)) == 0)
	})
	It("should be a slice with 1 items", func() {
		retSlice3 := SliceInSlice(s1, s4)
		Expect(len(retSlice3.([]int)) == 1)
		Expect(retSlice3.([]int)[0] == 4)
	})
})
