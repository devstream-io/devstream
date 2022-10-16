package configmanager

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadConfig", func() {
	var (
		configFilePath string
		config         *Config
		err            error
	)
	const (
		quickStartConfigPath = "../../../examples/quickstart.yaml"
		invalidConfigPath    = "not_exist.yaml"
	)

	JustBeforeEach(func() {
		config, err = NewManager(configFilePath).LoadConfig()
	})

	Context("when config file is valid", func() {
		BeforeEach(func() {
			configFilePath = quickStartConfigPath
		})

		Specify("the error should be nil", func() {
			Expect(config).NotTo(BeNil())
			Expect(err).To(BeNil())
		})
		It("should state filed correctly", func() {
			Expect(config.State.Backend).To(Or(Equal("local"), Equal("s3")))
			Expect(config.State.Options.StateFile).To(Equal("devstream.state"))
		})
		Specify("Tools Options cannot be empty", func() {
			Expect(len(config.Tools)).Should(BeNumerically(">", 0))
			Expect(config.Tools[0].Name).ShouldNot(BeEmpty())
			Expect(config.Tools[0].InstanceID).ShouldNot(BeEmpty())
			Expect(config.Tools[0].Options).ShouldNot(BeEmpty())
			Expect(len(config.Tools[0].DependsOn)).ShouldNot(BeNil())
		})
	})

	Context("when config file is invalid", func() {
		BeforeEach(func() {
			configFilePath = invalidConfigPath
		})
		It("should error", func() {
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		})
	})
})

var _ = Describe("config renders", func() {
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

	Describe("renderConfigs func", func() {
		When("config is wrong", func() {
			It("should return error if data is not valid yaml string", func() {
				notValidStr := "this is not valid yaml"
				notValidBytes := []byte(notValidStr)
				config, err := NewManager("").renderConfigs(notValidBytes, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(config).Should(BeNil())
			})

			It("should return error if core config is not valid", func() {
				config, err := NewManager("").renderConfigs(notValidCoreBytes, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(config).Should(BeNil())
			})

		})

		When("config is valid", func() {
			It("should generate config", func() {
				config, err := NewManager("").renderConfigs(validCoreBytes, emptyVariable, validToolBytes)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(config).ShouldNot(BeNil())
			})
		})
	})

	Describe("renderToolsFromCoreConfigAndConfigBytes func", func() {
		When("config error", func() {
			It("should return error if tool file is empty and tools config file is empty", func() {
				config := CoreConfig{}
				tools, err := NewManager("").renderToolsFromCoreConfigAndConfigBytes(&config, emptyVariable, emptyVariable)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("tools config is empty"))
				Expect(tools).Should(BeEmpty())

			})
		})
		When("tool config valid", func() {
			It("should generate tools array", func() {
				config := CoreConfig{}
				tools, err := NewManager("").renderToolsFromCoreConfigAndConfigBytes(&config, validToolBytes, emptyVariable)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(tools).ShouldNot(BeEmpty())
			})
		})
	})

	Describe("loadOriginalConfigFile func", func() {
		When("file name error", func() {
			It("should return error", func() {
				errorFileName := "error_file_name.yaml"
				_, err := NewManager(errorFileName).loadOriginalConfigFile()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})

	Describe("splitConfigFileBytes func", func() {
		When("input text not valid", func() {
			It("should return error if fomat not valid", func() {
				notValidContent := []byte(`
				---
				---
				---
				`)
				_, _, _, err := NewManager("").splitConfigFileBytes(notValidContent)
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
				_, _, _, err := NewManager("").splitConfigFileBytes(notValidContent)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("multiple sections"))
			})
		})
		When("input text is valid", func() {

			It("should return config bytes", func() {
				validArray := [][]byte{
					[]byte("---"),
					validCoreBytes,
					[]byte("---"),
					validToolBytes,
				}
				validInput := bytes.Join(validArray, []byte("\n"))
				coreConfigBytes, _, toolConfigBytes, err := NewManager("").splitConfigFileBytes(validInput)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(coreConfigBytes).ShouldNot(BeEmpty())
				Expect(toolConfigBytes).ShouldNot(BeEmpty())
			})
		})
	})

	Describe("checkConfigType func", func() {
		When("input is invalid", func() {
			It("should return false if config type not valid", func() {
				notValidType := []byte(`
test:
  - name: repo-scaffolding
    instanceID: default`)
				isValid, err := NewManager("").checkConfigType(notValidType, "core")
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(isValid).Should(BeFalse())
			})

			It("should return error if data is not valid yaml string", func() {
				notValidStr := "this is not valid yaml"
				notValidBytes := []byte(notValidStr)
				isValid, err := NewManager("").checkConfigType(notValidBytes, "core")
				Expect(err).Error().Should(HaveOccurred())
				Expect(isValid).Should(BeFalse())
			})
		})
		When("input is right", func() {
			It("should return true and error is nil", func() {
				isValid, err := NewManager("").checkConfigType(validCoreBytes, "core")
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(isValid).Should(BeTrue())

			})
		})
	})
})
