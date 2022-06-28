package jenkins_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/jenkins"
)

var _ = Describe("Validate", func() {

	Describe("test StorageClass name replacement", func() {
		Context("when StorageClass name is provided", func() {
			It("StorageClass should be replaced successfully", func() {
				valuesYaml := `
persistence:
            storageClass: custom-storage-class
          serviceAccount:
            create: false
            name: jenkins
          # nodePort, 30000 - 32767
          nodePort: 32000
          annotations: {}
`
				valuesYaml, err := jenkins.ReplaceStorageClass(valuesYaml)

				Expect(err).To(BeNil())
				Expect(valuesYaml).To(Equal(`
persistence:
            storageClass: jenkins-pv
          serviceAccount:
            create: false
            name: jenkins
          # nodePort, 30000 - 32767
          nodePort: 32000
          annotations: {}
`))
			})
		})
	})

	Context("when StorageClass name is not provided", func() {
		valuesYaml := `
persistence:
          serviceAccount:
            create: false
            name: jenkins
          # nodePort, 30000 - 32767
          nodePort: 32000
          annotations: {}
`
		_, err := jenkins.ReplaceStorageClass(valuesYaml)
		Expect(err).NotTo(BeNil())
	})
})
