package slicez

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slicez", func() {
	Context("Slice with type []int", func() {
		s1 := []int{1, 2, 3, 4}
		s2 := []int{1, 2, 3, 4}
		s3 := []int{1, 2, 3, 4, 5}
		s4 := []int{1, 2, 3}

		It("should be a slice with 0 items", func() {
			retSlice1 := SliceInSliceInt(s1, s2)
			Expect(len(retSlice1)).To(Equal(0))
			retSlice2 := SliceInSliceInt(s1, s3)
			Expect(len(retSlice2)).To(Equal(0))
		})
		It("should be a slice with 1 items", func() {
			retSlice3 := SliceInSliceInt(s1, s4)
			Expect(len(retSlice3)).To(Equal(1))
			Expect(retSlice3[0]).To(Equal(4))
		})
	})
	Context("Slice with type []string", func() {
		s11 := []string{"1", " 2", "3", "4"}
		s12 := []string{"1", " 2", "3", "4"}
		s13 := []string{"1", " 2", "3", "4", "5"}
		s14 := []string{"1", " 2", "3"}

		It("should be a slice with 0 items", func() {
			retSlice11 := SliceInSliceStr(s11, s12)
			Expect(len(retSlice11)).To(Equal(0))
			retSlice12 := SliceInSliceStr(s11, s13)
			Expect(len(retSlice12)).To(Equal(0))
		})
		It("should be a slice with 1 items", func() {
			retSlice13 := SliceInSliceStr(s11, s14)
			Expect(len(retSlice13)).To(Equal(1))
			Expect(retSlice13[0]).To(Equal("4"))
		})
	})
})
