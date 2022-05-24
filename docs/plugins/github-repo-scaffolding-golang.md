# github-repo-scaffolding-golang Plugin

This plugin bootstraps a GitHub repo with scaffolding code for a Golang web application.

_This plugin depends on the following environment variable:_

- GITHUB_TOKEN

Set it before using this plugin.

If you don't know how to create this token, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

*Tips:*

*1. If you run `dtm delete`, the repo on GitHub will be completely removed.*

*2. If the `Update` interface is called, the repo on github will be completely removed and recreated. However, given our current implementation, this interface shall not be called, as of in v0.2.0.*

## Usage

**Please note that the `owner` parameter is case-sensitive.**

```yaml
tools:
# name of the tool
- name: github-repo-scaffolding-golang
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
    # the repo which you'd like to create; please change the value below.
    repo: YOUR_REPO_NAME
    # the branch of the repo you'd like to hold the code
    branch: main
    # the image repo you'd like to push the container image; please change the value below.
    image_repo: YOUR_DOCKERHUB_USERNAME/YOUR_DOCKERHUB_REPOSITORY
```

Replace the following from the config above:

- `YOUR_GITHUB_USERNAME`
- `YOUR_ORGANIZATION_NAME`
- `YOUR_REPO_NAME`
- `YOUR_DOCKERHUB_USERNAME`
- `YOUR_DOCKERHUB_REPOSITORY`

The "branch" in the example above is "main", but you can adjust accordingly.

Currently, all the parameters in the example above are mandatory.

## Outputs

This plugin has three outputs:

- `owner`
- `repo`
- `repoURL` (example: "https://github.com/IronCore864/test.git")

If, for example, you want to use the outputs as inputs for another plugin, you can refer to the following example:

```yaml
---
tools:
- name: go-webapp-repo
  plugin: github-repo-scaffolding-golang
  options:
    owner: YOUR_GITHUB_USERNAME
    repo: YOUR_REPO_NAME
    branch: main
    image_repo: YOUR_DOCKERHUB_REPOSITORY
- name: golang-demo-actions
  plugin: githubactions-golang
  dependsOn: ["go-webapp-repo.github-repo-scaffolding-golang"]
  options:
    owner: ${{go-webapp-repo.github-repo-scaffolding-golang.outputs.owner}}
    repo: ${{go-webapp-repo.github-repo-scaffolding-golang.outputs.repo}}
    language:
      name: go
      version: "1.17"
    branch: main
    build:
      enable: True
    test:
      enable: True
      coverage:
        enable: True
    docker:
      enable: False
```

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.PLUGIN.outputs.var}}` is the syntax for using an output.
