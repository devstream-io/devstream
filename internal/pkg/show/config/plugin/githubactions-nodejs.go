package plugin

var GithubActionsNodejsDefaultConfig = `tools:
# name of the tool
- name: githubactions-nodejs
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    # the repo where you'd like to setup GitHub Actions; please change the value below to an existing repo.
    repo: YOUR_REPO_NAME
    # programming language specific settings
    language:
      name: nodejs
      # version of the language
      version: "16.14"
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main`
