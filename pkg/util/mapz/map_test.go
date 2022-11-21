package mapz_test

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/mapz"
)

var _ = Describe("Mapz", func() {
	keys := []string{
		"key1", "key2",
	}
	value := fmt.Errorf("error")
	expectMap := map[string]error{
		"key1": value,
		"key2": value,
	}
	retMap1 := mapz.FillMapWithStrAndError(keys, value)
	It("should be a map with 2 items", func() {
		Expect(retMap1).Should(Equal(expectMap))
	})

	retMap2 := mapz.FillMapWithStrAndError(nil, value)
	It("should be a map with 0 item", func() {
		Expect(len(retMap2)).To(Equal(0))
	})
})

func BenchmarkFillMapWithStrAndError(b *testing.B) {
	keys_length := 100
	keys := make([]string, 0, keys_length)
	for i := 0; i < keys_length; i++ {
		keys = append(keys, "key"+strconv.Itoa(i))
	}
	value := fmt.Errorf("error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapz.FillMapWithStrAndError(keys, value)
	}
	b.StopTimer()
}

var _ = Describe("Merge func", func() {
	var (
		src1, src2 map[string]any
	)
	BeforeEach(func() {
		src1 = map[string]any{
			"test1": "test1",
			"test":  "test",
		}
		src2 = map[string]any{
			"test1": "cover",
			"test2": "test2",
		}
	})
	It("should merge maps", func() {
		allMaps := mapz.Merge(src1, src2)
		Expect(allMaps).Should(Equal(map[string]any{
			"test1": "cover",
			"test":  "test",
			"test2": "test2",
		}))
	})
})
