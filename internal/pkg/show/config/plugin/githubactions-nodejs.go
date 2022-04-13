package plugin

var GithubActionsNodejsDefaultConfig = `tools:
- name: nodejs-demo-app
  # name of the plugin
  plugin: githubactions-nodejs
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
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
