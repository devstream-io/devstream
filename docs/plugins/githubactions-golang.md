# githubactions-golang Plugin

This plugin creates some Golang GitHub Actions workflows.

## Usage

_This plugin depends on the following environment variable:_

- GITHUB_TOKEN

Set it before using this plugin.

If you don't know how to create this token, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

_If Docker image build/push is enabled (see the example below), you also need to set the following two environment variables:_
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN

```yaml
tools:
# name of the tool
- name: githubactions-golang
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
      name: go
      # version of the language
      version: "1.17"
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
    build:
      # default to false
      enable: True
      # build command, OPTIONAL, the given value below is default value
      command: "go build ./..."
    test:
      # default to false
      enable: True
      # test command, OPTIONAL, the given value below is default value
      command: "go test ./..."
      coverage:
        # default to false
        enable: False
        # go test profile subcommand, OPTIONAL, the given value below is default value
        profile: "-race -covermode=atomic"
        output: "coverage.out"
    docker:
      # docker build/push related, default to false
      enable: True
      registry:
        # dockerhub or harbor, default to dockerhub
        type: dockerhub
        # dockerhub/harbor username
        username: YOUR_DOCKERHUB_USERNAME
        # dockerhub/harbor image repository name
        repository: YOUR_DOCKERHUB_REPOSITORY
```

Some parameters are optional. See the default values and optional parameters in the example above.

## Use Together with the `github-repo-scaffolding-golang` Plugin

This plugin can be used together with the `github-repo-scaffolding-golang` plugin (see document [here](./github-repo-scaffolding-golang.md).)

For example, you can first use `github-repo-scaffolding-golang` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `github-repo-scaffolding-golang` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `github-repo-scaffolding-golang`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: IronCore864
    repo: go-webapp-devstream-demo
    branch: main
    image_repo: ironcore864/go-webapp-devstream-demo
- name: githubactions-golang
  instanceID: default
  dependsOn: ["github-repo-scaffolding-golang.default"]
  options:
    owner: ${{github-repo-scaffolding-golang.default.outputs.owner}}
    repo: ${{github-repo-scaffolding-golang.default.outputs.repo}}
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
      enable: True
      registry:
        type: dockerhub
        username: [[ dockerhubUsername ]]
        repository: ${{github-repo-scaffolding-golang.default.outputs.repo}}
```

In the example above:

- We put `ggithub-repo-scaffolding-golang.default` as dependency by using the `dependsOn` keyword.
- We used `github-repo-scaffolding-golang.default`'s output as input for the `githubactions-golang` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.var}}` is the syntax for using an output.
