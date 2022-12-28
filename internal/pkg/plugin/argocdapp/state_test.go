package argocdapp

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

var _ = Describe("getStaticStatus func", func() {
	var option configmanager.RawOptions
	BeforeEach(func() {
		option = configmanager.RawOptions{
			"app": map[string]any{
				"name":      "test",
				"namespace": "test_namespace",
			},
			"source": map[string]any{
				"repoURL":   "test.com",
				"path":      "test_path",
				"valueFile": "val_file",
			},
			"destination": map[string]any{
				"server":    "test_server",
				"namespace": "dst_namespace",
			},
		}
	})
	It("should return status", func() {
		state, err := getStaticStatus(option)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(state).Should(Equal(statemanager.ResourceStatus{
			"app": map[string]any{
				"namespace": "test_namespace",
				"name":      "test",
			},
			"src": map[string]any{
				"repoURL":   "test.com",
				"path":      "test_path",
				"valueFile": "val_file",
			},
			"dest": map[string]any{
				"server":    "test_server",
				"namespace": "dst_namespace",
			},
		}))
	})
})
