package plugin

var JiraGithubDefaultConfig = `tools:
# name of the tool
- name: jira-github-integ
  # id of the tool instance
  instanceID: default
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
  # options for the plugin
  options:
    # the repo's owner
    owner: YOUR_GITHUB_USERNAME
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    # the repo where you'd like to setup GitHub Actions
    repo: YOUR_REPO_NAME
    # "base url: https://id.atlassian.net"
    jiraBaseUrl: https://JIRA_ID.atlassian.net
    # "need real user email in cloud Jira"
    jiraUserEmail: JIRA_USER_EMAIL
    # "get it from project url, like 'HEAP' from https://merico.atlassian.net/jira/software/projects/HEAP/pages"
    jiraProjectKey: JIRA_PROJECT_KEY 
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main`
