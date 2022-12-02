[[ if eq .RepoType "gitlab" ]]
// gitlab modules
import com.dabsquared.gitlabjenkins.GitLabPushTrigger
import com.dabsquared.gitlabjenkins.trigger.filter.BranchFilterType;
import com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty;
import org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition;
import hudson.plugins.git.GitSCM;
import hudson.plugins.git.BranchSpec;
import hudson.plugins.git.SubmoduleConfig;
import hudson.plugins.git.UserRemoteConfig
import hudson.plugins.git.extensions.impl.GitLFSPull;
import org.jenkinsci.plugins.workflow.flow.FlowDefinition
import static com.google.common.collect.Lists.newArrayList;
[[ end ]]

[[ if eq .RepoType "github" ]]
// github modules
import org.jenkinsci.plugins.github_branch_source.GitHubSCMSource
import jenkins.branch.BranchSource
import hudson.util.PersistedList
[[ end ]]

// common modules
import jenkins.model.Jenkins;
import org.jenkinsci.plugins.workflow.job.WorkflowJob;
import org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject;
import com.cloudbees.hudson.plugins.folder.*


// --> 1. init variables
Jenkins jenkinsInstance = Jenkins.instance
String repoType = "[[ .RepoType ]]"
String jobName = "[[ .JobName ]]";
String folderName = "[[ .FolderName ]]"
String jenkinsFileName = "Jenkinsfile"
String repoCredentialId = "[[ .RepoCredentialsId ]]"
// --> 2. init jobRef
Object jobRef = createJobWithFolder(jenkinsInstance, jobName, folderName, repoType)
// set job display name
jobRef.setDisplayName("[[ .JobName ]]")



[[ if eq .RepoType "gitlab" ]]
// --> 3. config gitlab related config
// config scm for Jenkinsfile
UserRemoteConfig userRemoteConfig = new UserRemoteConfig("[[ .RepositoryURL ]]", jobName, null, repoCredentialId)

branches = newArrayList(new BranchSpec("*/[[ .Branch ]]"))
doGenerateSubmoduleConfigurations = false
submoduleCfg = null
browser = null
gitTool = null
extensions = []
GitSCM scm = new GitSCM([userRemoteConfig], branches, doGenerateSubmoduleConfigurations, submoduleCfg, browser, gitTool, extensions)

FlowDefinition flowDefinition = (FlowDefinition) new CpsScmFlowDefinition(scm, jenkinsFileName)
jobRef.setDefinition(flowDefinition)

// config gitlab trigger
def gitlabTrigger = new GitLabPushTrigger()
gitlabTrigger.setSecretToken("[[ .SecretToken ]]")
gitlabTrigger.setTriggerOnPush(true)
gitlabTrigger.setTriggerOnMergeRequest(true)
gitlabTrigger.setBranchFilterType(BranchFilterType.RegexBasedFilter)
gitlabTrigger.setSourceBranchRegex(".*")
gitlabTrigger.setTargetBranchRegex("[[ .Branch ]]")

jobRef.addTrigger(gitlabTrigger)
def gitlabConnection = new GitLabConnectionProperty("[[ .GitlabConnection ]]")
jobRef.addProperty(gitlabConnection)
[[ end ]]

[[ if eq .RepoType "github" ]]
// --> 3. config github related config
jobRef.getProjectFactory().setScriptPath(jenkinsFileName)
GitHubSCMSource githubSource = new GitHubSCMSource("[[ .RepoOwner ]]", "[[ .RepoName ]]", "[[ .RepoURL ]]", true)
githubSource.setCredentialsId(repoCredentialId)
githubSource.setBuildOriginBranch(true)
githubSource.setBuildOriginPRMerge(true)
githubSource.setBuildForkPRMerge(false)
BranchSource branchSource = new BranchSource(githubSource)
PersistedList sources = jobRef.getSourcesList()
sources.clear()
sources.add(branchSource)
[[ end ]]

// create job
jobRef.save()

def createJobWithFolder(Jenkins jenkins, String jobName, String folderName, String repoType) {
    Object folderPath = null
    if (folderName != "") {
        def folder = jenkins.getItem(folderName)
        // create folder if it not exist
        if (folder == null) {
          folder = jenkins.createProject(Folder.class, folderName)
        }
        folderPath = folder
    } else {
        folderPath = jenkins
    }
    job = folderPath.getItem(jobName)
    if (job != null) {
        return job
    }
    def jobType = WorkflowJob
    if (repoType == "github") {
        jobType = WorkflowMultiBranchProject
    }
    return folderPath.createProject(jobType, jobName)
}
