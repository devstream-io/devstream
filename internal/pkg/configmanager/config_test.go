package configmanager

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		emptyVariable     []byte
		validToolBytes    []byte
		validCoreBytes    []byte
		notValidCoreBytes []byte
	)

	BeforeEach(func() {
		validToolBytes = []byte(`
tools:
  - name: repo-scaffolding
    instanceID: default
    options:
      owner: YOUR_GITHUB_USERNAME_CASE_SENSITIVE
      repo: go-webapp-devstream-demo
      branch: main
      image_repo: YOUR_DOCKER_USERNAME/go-webapp-devstream-demo`)
		notValidCoreBytes = []byte(`
varFile: tests
toolFile: test
state:
  backend: not_exist_backends`)
		validCoreBytes = []byte(`
varFile: ""
toolFile: ""
state:
  backend: local
  options:
    stateFile: devstream.state`)
	})
	Describe("LoadConfig yaml", func() {
		configStateObj, err := LoadConfig("../../../examples/quickstart.yaml")
		Context("when the Yaml parses successfully", func() {
			It("should state filed correctly", func() {
				Expect(configStateObj.State.Backend).To(Or(Equal("local"), Equal("s3")))
				Expect(configStateObj.State.Options.StateFile).To(Equal("devstream.state"))
			})
			Specify("Tools Options cannot be empty", func() {
				Expect(len(configStateObj.Tools)).ShouldNot(BeNil())
				Expect(configStateObj.Tools[0].Name).ShouldNot(BeEmpty())
				Expect(configStateObj.Tools[0].InstanceID).ShouldNot(BeEmpty())
				Expect(configStateObj.Tools[0].Options).ShouldNot(BeEmpty())
				Expect(len(configStateObj.Tools[0].DependsOn)).ShouldNot(BeNil())
			})
		})
		Context("when the Yaml parses fails", func() {
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("renderConfigs func", func() {
		Context("when config wrong", func() {
			It("should return error if data is not valid yaml string", func() {
				notValidStr := "this is not valid yaml"
				notValidBytes := []byte(notValidStr)
				config, err := renderConfigs(notValidBytes, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(config).Should(BeNil())
			})

			It("should return error if core config is not valid", func() {
				config, err := renderConfigs(notValidCoreBytes, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(config).Should(BeNil())
			})

		})

		Context("when config all valid", func() {
			It("should generate config", func() {
				config, err := renderConfigs(validCoreBytes, emptyVariable, validToolBytes)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(config).ShouldNot(BeNil())
			})
		})
	})

	Describe("renderToolsFromCoreConfigAndConfigBytes func", func() {
		Context("when config error", func() {
			It("should return error if tool file is empty and tools config file is empty", func() {
				config := CoreConfig{}
				tools, err := renderToolsFromCoreConfigAndConfigBytes(&config, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("tools config is empty"))
				Expect(tools).Should(BeEmpty())

			})
		})
		Context("when tool config valid", func() {
			It("should generate tools array", func() {
				config := CoreConfig{}
				tools, err := renderToolsFromCoreConfigAndConfigBytes(&config, validToolBytes, emptyVariable)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(tools).ShouldNot(BeEmpty())
			})
		})
	})

	Describe("loadOriginalConfigFile func", func() {
		Context("when file name error", func() {
			It("should return error", func() {
				errorFileName := "error_file_name.yaml"
				_, err := loadOriginalConfigFile(errorFileName)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})

	Describe("SplitConfigFileBytes func", func() {
		Context("when input text not valid", func() {
			It("should return error if fomat not valid", func() {
				notValidContent := []byte(`
				---
				---
				---
				`)
				_, _, _, err := SplitConfigFileBytes(notValidContent)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("invalid config format"))
			})
			It("should return error if duplicate config", func() {
				notValidContent := []byte(`
---
tools:
  - name: repo-scaffolding
    instanceID: default
---
tools:
  - name: repo-scaffolding
    instanceID: default`)
				_, _, _, err := SplitConfigFileBytes(notValidContent)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("multiple sections"))
			})
		})
		Context("when input text is valid", func() {

			It("should return config bytes", func() {
				validArray := [][]byte{
					[]byte("---"),
					validCoreBytes,
					[]byte("---"),
					validToolBytes,
				}
				validInput := bytes.Join(validArray, []byte("\n"))
				coreConfigBytes, _, toolConfigBytes, err := SplitConfigFileBytes(validInput)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(coreConfigBytes).ShouldNot(BeEmpty())
				Expect(toolConfigBytes).ShouldNot(BeEmpty())
			})
		})
	})

	Describe("checkConfigType func", func() {
		Context("when input not valid", func() {
			It("should return false if config type not valid", func() {
				notValidType := []byte(`
test:
  - name: repo-scaffolding
    instanceID: default`)
				isValid, err := checkConfigType(notValidType, "core")
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(isValid).Should(BeFalse())
			})

			It("should return error if data is not valid yaml string", func() {
				notValidStr := "this is not valid yaml"
				notValidBytes := []byte(notValidStr)
				isValid, err := checkConfigType(notValidBytes, "core")
				Expect(err).Error().Should(HaveOccurred())
				Expect(isValid).Should(BeFalse())
			})
		})
		Context("when input is right", func() {
			It("should return true and error is nil", func() {
				isValid, err := checkConfigType(validCoreBytes, "core")
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(isValid).Should(BeTrue())

			})
		})
	})
})
