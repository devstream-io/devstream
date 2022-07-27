package jenkinsgithub

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("All", func() {
	Describe("renderGitHubInteg func", func() {
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

	Describe("renderJobXml func", func() {
		It("should return the correct job xml", func() {
			opts := &JobOptions{
				JobName:              "job-pr",
				PrPipelineScriptPath: "Jenkinsfile-pr",
				GitHubRepoURL:        "https://github.com/aFlyBird0/jenkins-file-test",
				AdminList:            []string{"aFlyBird0", "Bird"},
				CredentialsID:        jenkinsCredentialID,
			}

			rendered, err := renderJobXml(jobPrTemplate, opts)
			Expect(err).To(BeNil())

			expected := `<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@1207.ve6191ff089f8">
    <description>A pipeline is only triggered when pr.</description>
    <keepDependencies>false</keepDependencies>
    <properties>
        <com.coravy.hudson.plugins.github.GithubProjectProperty plugin="github@1.34.4">
            <projectUrl>https://github.com/aFlyBird0/jenkins-file-test</projectUrl>
            <displayName></displayName>
        </com.coravy.hudson.plugins.github.GithubProjectProperty>
        <org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
            <triggers>
                <org.jenkinsci.plugins.ghprb.GhprbTrigger plugin="ghprb@1.42.2">
                    <spec>H/5 * * * *</spec>
                    <configVersion>3</configVersion>
                    <adminlist>aFlyBird0 Bird </adminlist>
                    <allowMembersOfWhitelistedOrgsAsAdmin>false</allowMembersOfWhitelistedOrgsAsAdmin>
                    <orgslist></orgslist>
                    <cron>H/5 * * * *</cron>
                    <buildDescTemplate></buildDescTemplate>
                    <onlyTriggerPhrase>false</onlyTriggerPhrase>
                    <useGitHubHooks>true</useGitHubHooks>
                    <permitAll>false</permitAll>
                    <whitelist></whitelist>
                    <autoCloseFailedPullRequests>false</autoCloseFailedPullRequests>
                    <displayBuildErrorsOnDownstreamBuilds>false</displayBuildErrorsOnDownstreamBuilds>
                    <whiteListTargetBranches>
                        <org.jenkinsci.plugins.ghprb.GhprbBranch>
                            <branch></branch>
                        </org.jenkinsci.plugins.ghprb.GhprbBranch>
                    </whiteListTargetBranches>
                    <blackListTargetBranches>
                        <org.jenkinsci.plugins.ghprb.GhprbBranch>
                            <branch></branch>
                        </org.jenkinsci.plugins.ghprb.GhprbBranch>
                    </blackListTargetBranches>
                    <gitHubAuthId>3a3b9ece-ad38-4209-8808-a37fbe74cc95</gitHubAuthId>
                    <triggerPhrase></triggerPhrase>
                    <skipBuildPhrase>.*\[skip\W+ci\].*</skipBuildPhrase>
                    <blackListCommitAuthor></blackListCommitAuthor>
                    <blackListLabels></blackListLabels>
                    <whiteListLabels></whiteListLabels>
                    <includedRegions></includedRegions>
                    <excludedRegions></excludedRegions>
                    <extensions>
                        <org.jenkinsci.plugins.ghprb.extensions.status.GhprbSimpleStatus>
                            <commitStatusContext></commitStatusContext>
                            <triggeredStatus></triggeredStatus>
                            <startedStatus></startedStatus>
                            <statusUrl></statusUrl>
                            <addTestResults>false</addTestResults>
                        </org.jenkinsci.plugins.ghprb.extensions.status.GhprbSimpleStatus>
                    </extensions>
                </org.jenkinsci.plugins.ghprb.GhprbTrigger>
            </triggers>
        </org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
    </properties>
    <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps@2759.v87459c4eea_ca_">
        <scm class="hudson.plugins.git.GitSCM" plugin="git@4.11.3">
            <configVersion>2</configVersion>
            <userRemoteConfigs>
                <hudson.plugins.git.UserRemoteConfig>
                    <name>origin</name>
                    <refspec>+refs/pull/*:refs/remotes/origin/pr/*</refspec>
                    <url>https://github.com/aFlyBird0/jenkins-file-test</url>
                    <credentialsId>credential-by-devstream-jenkins-github-integ</credentialsId>
                </hudson.plugins.git.UserRemoteConfig>
            </userRemoteConfigs>
            <branches>
                <hudson.plugins.git.BranchSpec>
                    <name>${sha1}</name>
                </hudson.plugins.git.BranchSpec>
            </branches>
            <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
            <submoduleCfg class="empty-list"/>
            <extensions/>
        </scm>
        <scriptPath>Jenkinsfile-pr</scriptPath>
        <lightweight>false</lightweight>
    </definition>
    <triggers/>
    <disabled>false</disabled>
</flow-definition>
`
			Expect(rendered).To(Equal(expected))
		})
	})
})
