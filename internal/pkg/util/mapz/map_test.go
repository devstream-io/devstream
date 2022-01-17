package mapz

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mapz", func() {
	keys := []string{
		"key1", "key2",
	}
	value := fmt.Errorf("error")
	retMap1 := FillMapWithStrAndError(keys, value)
	It("should be a map with 2 items", func() {
		Expect(len(retMap1)).To(Equal(2))
		v1, ok := retMap1["key1"]
		Expect(ok).To(Equal(true))
		Expect(v1).To(Equal(value))
		v2, ok := retMap1["key2"]
		Expect(ok).To(Equal(true))
		Expect(v2).To(Equal(value))
	})

	retMap2 := FillMapWithStrAndError(nil, value)
	It("should be a map with 0 item", func() {
		Expect(len(retMap2)).To(Equal(0))
	})
})
