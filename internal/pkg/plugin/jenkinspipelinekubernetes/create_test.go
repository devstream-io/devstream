package jenkinspipelinekubernetes

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	Describe("renderJobXml function", func() {
		var opts *JobXmlOptions
		BeforeEach(func() {
			opts = &JobXmlOptions{
				GitHubRepoURL:      "https://github.com/xxx/jenkins-file-test.git",
				CredentialsID:      "credential-jenkins-pipeline-kubernetes-by-devstream",
				PipelineScriptPath: "this-is-pipeline-script-path",
			}
		})

		It("should return the correct xml", func() {
			xml := renderJobXml(jobTemplate, opts)
			expect := `<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@1189.va_d37a_e9e4eda_">
    <description></description>
    <keepDependencies>false</keepDependencies>
    <properties/>
    <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps@2725.v7b_c717eb_12ce">
        <scm class="hudson.plugins.git.GitSCM" plugin="git@4.11.3">
            <configVersion>2</configVersion>
            <userRemoteConfigs>
                <hudson.plugins.git.UserRemoteConfig>
                    <url>https://github.com/xxx/jenkins-file-test.git</url>
                    <credentialsId>credential-jenkins-pipeline-kubernetes-by-devstream</credentialsId>
                </hudson.plugins.git.UserRemoteConfig>
            </userRemoteConfigs>
            <branches>
                <hudson.plugins.git.BranchSpec>
                    <name>*/main</name>
                </hudson.plugins.git.BranchSpec>
            </branches>
            <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
            <submoduleCfg class="empty-list"/>
            <extensions/>
        </scm>
        <scriptPath>this-is-pipeline-script-path</scriptPath>
        <lightweight>false</lightweight>
    </definition>
    <triggers/>
    <disabled>false</disabled>
</flow-definition>
`
			Expect(xml).To(Equal(expect))
		})
	})
})
