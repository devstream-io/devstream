package k8s

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Configmap", func() {
	It("should work", func() {
		data := map[string]any{
			"a": "A",
			"b": map[string]string{
				"c": "C",
			},
		}

		b, err := NewBackend("devstream", "state")
		Expect(err).To(BeNil())
		Expect(b).ToNot(BeNil())

		dataBytes, err := yaml.Marshal(data)
		Expect(err).To(BeNil())
		dataStr := string(dataBytes)

		// test applyConfigMap
		cm, err := b.applyConfigMap(dataStr)
		Expect(err).To(BeNil())
		Expect(cm.Data[stateKey]).To(Equal(dataStr))

		// test getOrCreateConfigMap
		cm, err = b.getOrCreateConfigMap()
		Expect(err).To(BeNil())
		Expect(cm).ToNot(BeNil())
		Expect(cm.Data[stateKey]).To(Equal(dataStr))
	})
})
