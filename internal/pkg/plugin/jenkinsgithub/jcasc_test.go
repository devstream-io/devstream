package jenkinsgithub

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderGitHubInteg func", func() {
	It("should return the correct github integration config", func() {
		opts := &GitHubIntegOptions{
			AdminList:          []string{"aFlyBird0", "Bird"},
			CredentialsID:      jenkinsCredentialID,
			GithubAuthID:       githubAuthID,
			JenkinsURLOverride: "https://891e-125-111-206-162.ap.ngrok.io/",
		}

		rendered, err := renderGitHubInteg(opts)
		Expect(err).To(BeNil())

		expected := `unclassified:
  ghprbTrigger:
    adminlist: "aFlyBird0 Bird "
    autoCloseFailedPullRequests: false
    cron: "H/5 * * * *"
    extensions:
    - ghprbSimpleStatus:
        addTestResults: false
        showMatrixStatus: false
    githubAuth:
    - credentialsId: "credential-by-devstream-jenkins-github-integ"
      description: "Anonymous connection"
      id: "3a3b9ece-ad38-4209-8808-a37fbe74cc95"
      jenkinsUrl: "https://891e-125-111-206-162.ap.ngrok.io/"
      serverAPIUrl: "https://api.github.com"
    manageWebhooks: true
    okToTestPhrase: ".*ok\\W+to\\W+test.*"
    requestForTestingPhrase: "Can one of the admins verify this patch?"
    retestPhrase: ".*test\\W+this\\W+please.*"
    skipBuildPhrase: ".*\\[skip\\W+ci\\].*"
    useComments: false
    useDetailedComments: false
    whitelistPhrase: ".*add\\W+to\\W+whitelist.*"
`
		Expect(rendered).To(Equal(expected))
	})
})
