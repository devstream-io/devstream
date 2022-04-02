package k8s_test

import (
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
	//"github.com/devstream-io/devstream/pkg/util/k8s"
)

var _ = Describe("K8S", func() {
	Context("Namespace", func() {
		// TODO(daniel-hutao): the code below is only used local now for the k8s test env not exist at GitHub.
		//c, err := k8s.NewClient()
		//Expect(err).NotTo(HaveOccurred())

		// create
		//err = c.CreateNamespace("test-ds-ns")
		//Expect(err).NotTo(HaveOccurred())

		// get
		//ns, err := c.GetNamespace("test-ds-ns")
		//Expect(err).NotTo(HaveOccurred())
		//Expect(ns).NotTo(Equal(nil))
		//Expect(ns.Name).To(Equal("test-ds-ns"))

		// delete
		//err = c.DeleteNamespace("monitoring")
		//Expect(err).NotTo(HaveOccurred())
	})
})
