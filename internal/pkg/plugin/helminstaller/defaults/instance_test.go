package defaults

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("instance test", func() {
	var (
		instanceID string
		instance   *HelmAppInstance
	)

	// main logic to get the default instance
	JustBeforeEach(func() {
		instance = GetDefaultHelmAppInstanceByInstanceID(instanceID)
	})

	Context("GetDefaultHelmAppInstanceByInstanceID", func() {
		When("instanceID is not exists", func() {
			BeforeEach(func() {
				instanceID = "not-exists"
			})
			It("should return nil", func() {
				// main logic is in JustBeforeEach
				Expect(instance).To(BeNil())
			})
		})

		When("instanceID is exists(argocd)", func() {
			JustAfterEach(func() {
				Expect(instance).NotTo(BeNil())
				Expect(instance.Name).To(Equal("argocd"))
				Expect(*instance.HelmOptions).To(Equal(DefaultConfigWithArgoCD))
				pointer1 := reflect.ValueOf(instance.StatusGetter).Pointer()
				pointer2 := reflect.ValueOf(GetArgoCDStatus).Pointer()
				Expect(pointer1).To(Equal(pointer2))
			})
			When(`instanceID is the same as "argocd"`, func() {
				BeforeEach(func() {
					instanceID = "argocd"
				})
				It("should return instance", func() {
					// main logic is in JustAfterEach
				})
			})

			When(`instanceID has prefix "argocd"`, func() {
				BeforeEach(func() {
					instanceID = "argocd-001"
				})
				It("should return instance", func() {
					// main logic is in JustAfterEach
				})
			})
		})
	})
})
