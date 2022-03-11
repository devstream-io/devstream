## 1 `jira-github-integ` Plugin

This plugin integrates Jira with your GitHub repo.

## 2 Usage:

_Please confirm the preconditions:_

- Jira language must be English
- There should be an existing Jira project

_This plugin depends on the following two environment variables:_

- JIRA_API_TOKEN
- GITHUB_TOKEN

Set the values accordingly before using this plugin.


To create a Jira API token, see [here](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/).

```yaml
tools:
- name: default
  plugin:
    # kind of this plugin
    kind: jira-github-integ
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.2.0
  # options for the plugin
  options:
    # the repo's owner
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions
    repo: YOUR_REPO_NAME
    # "base url: https://id.atlassian.net"
    jiraBaseUrl: https://JIRA_ID.atlassian.net
    # "need real user email in cloud Jira"
    jiraUserEmail: JIRA_USER_EMAIL
    # "get it from project url, like 'HEAP' from https://merico.atlassian.net/jira/software/projects/HEAP/pages"
    jiraProjectKey: JIRA_PROJECT_KEY 
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

Currently, all the parameters in the example above are mandatory.
