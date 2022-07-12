package gitlabcedocker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options", func() {

	var opts Options

	BeforeEach(func() {
		opts = Options{
			GitLabHome:        "/srv/gitlab",
			Hostname:          "gitlab.example.com",
			SSHPort:           22,
			HTTPPort:          80,
			HTTPSPort:         443,
			RmDataAfterDelete: false,
			ImageTag:          "rc",
		}
	})

	Describe("getVolumesDirFromOptions func", func() {
		Context("when the options is valid", func() {
			It("should return the volumes' directory", func() {
				volumesDirFromOptions := getVolumesDirFromOptions(opts)
				Expect(volumesDirFromOptions).To(Equal([]string{
					"/srv/gitlab/config",
					"/srv/gitlab/data",
					"/srv/gitlab/logs",
				}))
			})
		})
	})
})
