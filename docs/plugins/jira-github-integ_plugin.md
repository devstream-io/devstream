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
- name: jira-github-integ-default
  # plugin profile
  plugin:
    # kind of this plugin
    kind: jira-github-integ
    # version of the plugin
    version: 0.3.0
  # options for the plugin
  # checkout the version from the GitHub releases
  options:
    # the repo's owner
    owner: lfbdev
    # the repo where you'd like to setup GitHub Actions
    repo: opendeveloper
    # "base url: https://id.atlassian.net"
    jiraBaseUrl: https://merico.atlassian.net
    # "need real user email in cloud Jira"
    jiraUserEmail: fangbao.li@merico.dev
    # "get it from project url, like 'HEAP' from https://merico.atlassian.net/jira/software/projects/HEAP/pages"
    jiraProjectKey: HEAP 
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: master
```

Currently, all the parameters in the example above are mandatory.
