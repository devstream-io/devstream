package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipelinetemplate", func() {
	var _ = Describe("mergeMaps", func() {
		When("two maps have same keys", func() {
			It("should overwrite the fisrt map", func() {
				map1 := map[string]any{
					"key1": "value1",
					"key2": "value2",
				}
				map2 := map[string]any{
					"key2": "value-of-map2",
					"key3": "value2",
				}
				map3 := mergeMaps(map1, map2)
				mapExpected := map[string]any{
					"key1": "value1",
					"key2": "value-of-map2",
					"key3": "value2",
				}
				Expect(map3).Should(Equal(mapExpected))
			})
		})
	})
})
