# githubactions-nodejs Plugin

This plugin creates Nodejs GitHub Actions workflows.

## Usage

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

If you don't know how to create this token, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

```yaml
tools:
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
    branch: main
```

All parameters are mandatory.
