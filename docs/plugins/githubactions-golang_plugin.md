## 1 `githubactions-golang` Plugin

This plugin creates some Golang GitHub Actions workflows.

## 2 Usage

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
- name: golang-demo-app
  plugin:
    # name of the plugin
    kind: githubactions-golang
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.3.0
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_KIND", "TOOL2_NAME.TOOL2_KIND" ]
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below to an existing repo.
    repo: YOURE_REPO_NAME
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
      # dockerhub image repo; please change the value below.
      repo: YOUR_DOCKERHUB_IMAGE_REPO_NAME
```

Some parameters are optional. See the default values and optional parameters in the example above.

## 3. Use Together with the `github-repo-scaffolding-golang` Plugin

This plugin can be used together with the `github-repo-scaffolding-golang` plugin (see document [here](./github-repo-scaffolding-golang_plugin.md).)

For example, you can first use `github-repo-scaffolding-golang` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `github-repo-scaffolding-golang` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `github-repo-scaffolding-golang`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
tools:
- name: go-webapp-repo
  plugin:
    kind: github-repo-scaffolding-golang
    version: 0.3.0
  options:
    owner: IronCore864
    repo: go-webapp-devstream-demo
    branch: main
    image_repo: ironcore864/go-webapp-devstream-demo
- name: golang-demo-actions
  plugin:
    kind: githubactions-golang
    version: 0.3.0
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

In the example above:

- We put `go-webapp-repo.github-repo-scaffolding-golang` as dependency by using the `dependsOn` keyword.
- We used `go-webapp-repo.github-repo-scaffolding-golang`'s output as input for the `githubactions-golang` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_KIND.outputs.var}}` is the syntax for using an output.
