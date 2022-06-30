package gitlabcedocker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabcedocker"
)

var _ = Describe("delete", func() {
	var options map[string]interface{}

	BeforeEach(func() {
		options = map[string]interface{}{
			"hostname":    "gitlab.devstream.io",
			"gitlab_home": "/srv/gitlab",
			"ssh_port":    8122,
			"http_port":   8180,
			"https_port":  8443,
		}
	})

	It("delete docker data successfully", func() {
		cmd, _ := gitlabcedocker.Delete(options)
		Expect(cmd).To(Equal(true))
	})
})
