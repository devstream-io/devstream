package gitlabcedocker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/docker"
)

var _ = Describe("Options", func() {

	var opts Options

	BeforeEach(func() {
		opts = Options{
			GitLabHome:        "/srv/gitlab",
			Hostname:          "gitlab.example.com",
			SSHPort:           8122,
			HTTPPort:          8180,
			HTTPSPort:         8443,
			RmDataAfterDelete: false,
			ImageTag:          "rc",
		}
	})

	Describe("getVolumesDirFromOptions func", func() {
		When("the options is valid", func() {
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

	Describe("buildDockerRunOptions func", func() {
		It("should build the docker run options successfully", func() {
			runOptsBuild := buildDockerRunOptions(opts)
			runOptsExpect := docker.RunOptions{
				ImageName:     "gitlab/gitlab-ce",
				ImageTag:      "rc",
				Hostname:      "gitlab.example.com",
				ContainerName: "gitlab",
				RestartAlways: true,
				PortPublishes: []docker.PortPublish{
					{HostPort: 8122, ContainerPort: 22},
					{HostPort: 8180, ContainerPort: 80},
					{HostPort: 8443, ContainerPort: 443},
				},
				Volumes: []docker.Volume{
					{HostPath: "/srv/gitlab/config", ContainerPath: "/etc/gitlab"},
					{HostPath: "/srv/gitlab/data", ContainerPath: "/var/opt/gitlab"},
					{HostPath: "/srv/gitlab/logs", ContainerPath: "/var/log/gitlab"},
				},
			}

			Expect(runOptsBuild).To(Equal(runOptsExpect))
		})

	})
})
