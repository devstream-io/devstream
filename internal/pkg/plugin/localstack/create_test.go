package localstack

import (
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Create", func() {

	It("should return the correct xml", func() {
		options := make(map[string]interface{})
		yaml_options := `
create_namespace: false
repo:
  name: localstack-charts
  url: https://localstack.github.io/helm-charts
chart:
  chart_name: localstack-charts/localstack
  release_name: localstack
  namespace: default
  wait: true
  timeout: 5m
  upgradeCRDs: true
  values_yaml: |
    debug: true
    updateStrategy:
    type: Recreate
`

		err := yaml.Unmarshal([]byte(yaml_options), &options)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		_, err = Create(options)
		Ω(err).ShouldNot(HaveOccurred())
	})

})
