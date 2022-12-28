package jenkins

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("jenkins basic method", func() {
	var (
		j              *jenkins
		url, namespace string
	)
	BeforeEach(func() {
		url = "testurl.com"
		namespace = "test_namespace"
		j = &jenkins{
			BasicInfo: &JenkinsConfigOption{
				URL:       url,
				Namespace: namespace,
			},
		}
	})
	It("should return basicInfo", func() {
		basicInfo := j.GetBasicInfo()
		Expect(basicInfo.URL).Should(Equal(url))
		Expect(basicInfo.Namespace).Should(Equal(namespace))
	})
})
