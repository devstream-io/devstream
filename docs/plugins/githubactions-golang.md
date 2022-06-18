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
--8<-- "githubactions-golang.yaml"
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
      version: "1.18"
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
