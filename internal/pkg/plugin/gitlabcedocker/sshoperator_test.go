package gitlabcedocker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabcedocker"
)

var _ = Describe("SSH operator", func() {
	var options gitlabcedocker.Options

	BeforeEach(func() {
		options = gitlabcedocker.Options{
			Hostname:   "gitlab.devstream.io",
			GitLabHome: "/srv/gitlab",
			SSHPort:    8122,
			HTTPPort:   8180,
			HTTPSPort:  8443,
		}
	})

	It("docker run command should be built correctly", func() {
		cmd := gitlabcedocker.BuildDockerRunCommand(options)
		expect := `
	docker run --detach \
	--hostname gitlab.devstream.io \
	--publish 8443:443 --publish 8180:80 --publish 8122:22 \
	--name gitlab \
	--restart always \
	--volume /srv/gitlab/config:/etc/gitlab \
	--volume /srv/gitlab/logs:/var/log/gitlab \
	--volume /srv/gitlab/data:/var/opt/gitlab \
	--shm-size 256m \
	gitlab/gitlab-ce:rc
	`
		Expect(cmd).To(Equal(expect))
	})
})
