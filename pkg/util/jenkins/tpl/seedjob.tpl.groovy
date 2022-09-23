[[ if eq .RepoType "gitlab" ]]
import com.dabsquared.gitlabjenkins.GitLabPushTrigger
import com.dabsquared.gitlabjenkins.trigger.filter.BranchFilterType;
import com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty;
[[ end ]]

import hudson.plugins.git.GitSCM;
import hudson.plugins.git.BranchSpec;
import hudson.plugins.git.SubmoduleConfig;
import hudson.plugins.git.UserRemoteConfig
import hudson.plugins.git.extensions.impl.GitLFSPull;

import hudson.util.Secret;
import jenkins.model.Jenkins;
import org.jenkinsci.plugins.workflow.job.WorkflowJob;
import org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition;
import org.jenkinsci.plugins.workflow.flow.FlowDefinition

import com.cloudbees.hudson.plugins.folder.*

import static com.google.common.collect.Lists.newArrayList;

// --> 1. this module is used to init jobRef
Object jobRef = null
def jobName = "[[ .JobName ]]";
Jenkins jenkins = Jenkins.instance
// Create job in folder if FolderName is config
[[ if .FolderName ]]
def folder = jenkins.getItem("[[ .FolderName ]]")
if (folder == null) {
  folder = jenkins.createProject(Folder.class, "[[ .FolderName ]]")
}
jobRef = folder.getItem(jobName)
if (jobRef == null) {
  oldJob = jenkins.getItem(jobName)
  if (oldJob.getClass() == WorkflowJob.class) {
    // Move any existing job into the folder
    Items.move(oldJob, folder)
  } else {
    // Create it if it doesn't
    jobRef = folder.createProject(WorkflowJob, jobName)
  }
}
[[ else ]]
// Create job directly
jobRef = jenkins.getItem(jobName)
if (jobRef == null) {
        jobRef = jenkins.createProject(WorkflowJob, jobName)
}
[[ end ]]

// set display name
jobRef.setDisplayName("[[ .JobName ]]")

// --> 2. this module is used to init jenkinsfile config
UserRemoteConfig userRemoteConfig = new UserRemoteConfig("[[ .RepositoryURL ]]", "[[ .JobName ]]", null, "[[ .RepoCredentialsId ]]")

branches = newArrayList(new BranchSpec("*/[[ .Branch ]]"))
doGenerateSubmoduleConfigurations = false
submoduleCfg = null
browser = null
gitTool = null
extensions = []
GitSCM scm = new GitSCM([userRemoteConfig], branches, doGenerateSubmoduleConfigurations, submoduleCfg, browser, gitTool, extensions)

FlowDefinition flowDefinition = (FlowDefinition) new CpsScmFlowDefinition(scm, "Jenkinsfile")
jobRef.setDefinition(flowDefinition)

// --> 3. config gitlab trigger
[[ if eq .RepoType "gitlab" ]]
def gitlabTrigger = new GitLabPushTrigger()
gitlabTrigger.setSecretToken("[[ .SecretToken ]]")
gitlabTrigger.setTriggerOnPush(true)
gitlabTrigger.setTriggerOnMergeRequest(true)
gitlabTrigger.setBranchFilterType(BranchFilterType.RegexBasedFilter)
gitlabTrigger.setSourceBranchRegex(".*")
gitlabTrigger.setTargetBranchRegex("master")

jobRef.addTrigger(gitlabTrigger)
def gitlabConnection = new GitLabConnectionProperty("[[ .GitlabConnection ]]")
jobRef.addProperty(gitlabConnection)
[[ end ]]

// create job
jobRef.save()
