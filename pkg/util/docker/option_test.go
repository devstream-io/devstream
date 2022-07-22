package docker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Option", func() {
	Describe("IfVolumesDiffer func", func() {
		var (
			volumes1            = []string{"/srv/gitlab/config", "/srv/gitlab/data", "/srv/gitlab/logs"}
			Volumes1ChangeOrder = []string{"/srv/gitlab/data", "/srv/gitlab/logs", "/srv/gitlab/config"}
			volumes1Missing     = []string{"/srv/gitlab/data", "/srv/gitlab/logs"}
			volumes2            = []string{"totally/different/path"}
		)

		var (
			volumesSrc, volumesDest []string
			differ                  bool
		)

		JustBeforeEach(func() {
			differ = IfVolumesDiffer(volumesSrc, volumesDest)
		})

		When("the volumes are the same but the order is changed", func() {
			BeforeEach(func() {
				volumesSrc = volumes1
				volumesDest = Volumes1ChangeOrder
			})

			It("should return false", func() {
				Expect(differ).To(BeFalse())
			})
		})

		When("the volumes are different(missing)", func() {
			BeforeEach(func() {
				volumesSrc = volumes1
				volumesDest = volumes1Missing
			})

			It("should return true", func() {
				Expect(differ).To(BeTrue())
			})
		})

		When("the volumes are totally different", func() {
			BeforeEach(func() {
				volumesSrc = volumes1
				volumesDest = volumes2
			})

			It("should return true", func() {
				Expect(differ).To(BeTrue())
			})
		})
	})
})
