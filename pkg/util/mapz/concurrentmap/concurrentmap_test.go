package concurrentmap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/mapz/concurrentmap"
)

var _ = Describe("Concurrentmap", func() {
	Context("CURD", func() {
		It("Should has 1 item", func() {
			cMap := concurrentmap.NewConcurrentMap("", "")
			cMap.Store("key1", "value1")
			v1, ok1 := cMap.Load("key1")
			Expect(ok1).To(BeTrue())
			Expect(v1.(string)).To(Equal("value1"))

			v2, ok2 := cMap.Load("key2")
			Expect(ok2).To(BeFalse())
			Expect(v2).To(BeNil())

			cMap.Delete("key1")
			v3, ok3 := cMap.Load("key1")
			Expect(ok3).To(BeFalse())
			Expect(v3).To(BeNil())
		})
	})
})
