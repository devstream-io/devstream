package configGetter_test

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/devstream-io/devstream/pkg/util/config/configGetter"
	"github.com/devstream-io/devstream/pkg/util/k8s"
)

var _ = Describe("Getter general", func() {

	var (
		key1      = "key1"
		key2      = "key2"
		key3      = "key3"
		keyNotSet = "keyNotSet"
	)

	var valueMap map[string]string = map[string]string{
		key1:      "value1",
		key2:      "value2",
		key3:      "value3",
		keyNotSet: "value4",
	}

	BeforeEach(func() {
		key1OriginalValue := os.Getenv(key1)
		key2OriginalValue := os.Getenv(key2)
		key3OriginalValue := os.Getenv(key3)

		os.Setenv(key1, valueMap[key1])
		os.Setenv(key2, valueMap[key2])
		os.Setenv(key3, valueMap[key3])

		DeferCleanup(func() {
			err := os.Setenv(key1, key1OriginalValue)
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(key2, key2OriginalValue)
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(key3, key3OriginalValue)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("get value from chain", func() {
		It("should return the first non-empty value", func() {
			getter := configGetter.NewItemGetterChain(
				configGetter.NewEnvGetter(key1),
				configGetter.NewEnvGetter(key2),
			)
			Expect(getter.Get()).To(Equal(valueMap[key1]))
		})

		It("should return error if all values are empty", func() {
			getter := configGetter.NewItemGetterChain(
				configGetter.NewEnvGetter(keyNotSet),
			)
			_, err := getter.Get()
			Expect(err).To(HaveOccurred())
		})

		It("should return right value if default value is set", func() {
			const defaultValue = "defaultValue"
			getter := configGetter.NewItemGetterChain(
				configGetter.NewEnvGetter(keyNotSet),
				configGetter.DefaultValue(defaultValue),
			)
			Expect(getter.Get()).To(Equal(defaultValue))
		})

	})

})

var _ = Describe("Each Getter", func() {
	Describe("EnvGetter", func() {
		// it's tested in the general test
	})

	Describe("Tool Options Getter", func() {
		var (
			rawOptions map[string]interface{} = map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": map[string]interface{}{
					"key4": "value4",
				},
				"key5": map[string]interface{}{
					"key6": []string{"value5", "value6"},
				},
			}
		)
		When("simple key", func() {
			It("should return the right value", func() {
				getter := configGetter.NewToolOptionsGetter("key1", rawOptions)
				Expect(getter.Get()).To(Equal("value1"))
			})
		})
		When("nested key", func() {
			It("should return the right value", func() {
				getter := configGetter.NewToolOptionsGetter("key3.key4", rawOptions)
				Expect(getter.Get()).To(Equal("value4"))
			})
		})

		When("key is not set", func() {
			It("should return empty string", func() {
				getter := configGetter.NewToolOptionsGetter("keyNotFound", rawOptions)
				Expect(getter.Get()).To(Equal(""))
			})
		})
	})

	Describe("viper getter", func() {
		// from viper official:
		// "Important: Viper configuration keys are case insensitive.
		// There are ongoing discussions about making that optional. "
		When("key is set", func() {
			When("key is lower case", func() {
				BeforeEach(func() {
					viper.Set("key1", "value1")
				})

				It("should return the right value", func() {
					getter := configGetter.NewViperGetter("KEY1")
					Expect(getter.Get()).To(Equal("value1"))
					getter = configGetter.NewViperGetter("key1")
					Expect(getter.Get()).To(Equal("value1"))
				})
			})

			When("key is upper case", func() {
				BeforeEach(func() {
					viper.Set("KEY2", "value2")
				})

				It("should return the right value", func() {
					getter := configGetter.NewViperGetter("KEY2")
					Expect(getter.Get()).To(Equal("value2"))
					getter = configGetter.NewViperGetter("key2")
					Expect(getter.Get()).To(Equal("value2"))
				})
			})

		})
		When("key is not set", func() {
			It("should return empty string", func() {
				getter := configGetter.NewViperGetter("keyNotFound")
				Expect(getter.Get()).To(Equal(""))
			})
		})

		When("key is set by os.SetEnv", func() {
			var key, value, keyUpper string

			BeforeEach(func() {
				// notice that the key is upper case
				key, value = "key_1", "value"
				keyUpper = strings.ToUpper(key)
				originValue := os.Getenv(keyUpper)
				os.Setenv(keyUpper, value)
				viper.AutomaticEnv()
				DeferCleanup(func() {
					err := os.Setenv(keyUpper, originValue)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			It("should return the right value", func() {
				getter := configGetter.NewViperGetter(key)
				Expect(getter.Get()).To(Equal(value))
			})
		})
	})

	Describe("K8s Secret Getter", func() {
		const (
			namespace, secretName = "test-ns", "test-secret"
			key1, value1          = "key1", "value1"
			keyNotSet             = "key-not-set"
		)
		BeforeEach(func() {
			// create a fake k8s client set set a secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
					Labels: map[string]string{
						"usage": "test",
					},
				},
				Data: map[string][]byte{
					key1: []byte(value1),
				},
			}
			fakeClient := fake.NewSimpleClientset(secret)
			k8s.UseFakeClient(fakeClient, nil)
		})

		When("key is set", func() {
			It("should return the right value", func() {
				getter := configGetter.NewK8sSecretGetter(key1, namespace, secretName)
				Expect(getter.Get()).To(Equal(value1))
			})
		})

		When("key is not set", func() {
			It("should return empty string", func() {
				getter := configGetter.NewK8sSecretGetter(keyNotSet, namespace, secretName)
				Expect(getter.Get()).To(Equal(""))
			})
		})
	})

	Describe("DefaultValue Getter", func() {
		// it's tested in the general test
	})
})
