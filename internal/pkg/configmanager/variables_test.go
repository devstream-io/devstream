package configmanager

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderConfigWithVariables", func() {
	When("render a config with variables", func() {
		It("should works fine", func() {
			var vars = map[string]interface{}{
				"foo1": "bar1",
				"foo2": "bar2",
			}
			result1, err1 := renderConfigWithVariables("[[ foo1 ]]/[[ foo2]]", vars)
			result2, err2 := renderConfigWithVariables("a[[ foo1 ]]/[[ foo2]]b", vars)
			result3, err3 := renderConfigWithVariables(" [[ foo1 ]] [[ foo2]] ", vars)
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(result1).To(Equal([]byte("bar1/bar2")))
			Expect(result2).To(Equal([]byte("abar1/bar2b")))
			Expect(result3).To(Equal([]byte(" bar1 bar2 ")))
		})
	})

	When("there are env variables", func() {
		const (
			envKey1, envKey2 = "GITHUB_TOKEN", "GITLAB_TOKEN"
			envVal1, envVal2 = "123456", "abcdef"
		)
		BeforeEach(func() {
			DeferCleanup(os.Setenv, envKey1, os.Getenv(envKey1))
			DeferCleanup(os.Setenv, envKey2, os.Getenv(envKey2))
		})

		It("should works fine", func() {
			err := os.Setenv(envKey1, envVal1)
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(envKey2, envVal2)
			Expect(err).NotTo(HaveOccurred())

			content := fmt.Sprintf(`
githubToken: [[env %s]]
gitlabToken: [[ env "%s"]]
`, envKey1, envKey2)
			expected := fmt.Sprintf(`
githubToken: %s
gitlabToken: %s
`, envVal1, envVal2)
			result, err := renderConfigWithVariables(content, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]byte(expected)))
		})
	})
})
