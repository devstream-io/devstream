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

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "githubactions-golang.yaml"
```

Some parameters are optional. See the default values and optional parameters in the example above.

## Use Together with the `repo-scaffolding` Plugin

This plugin can be used together with the `repo-scaffolding` plugin (see document [here](./repo-scaffolding.md).)

For example, you can first use `repo-scaffolding` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `repo-scaffolding` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `repo-scaffolding`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
tools:
- name: repo-scaffolding
  instanceID: golang-github
  options:
    owner: IronCore864
    repo: go-webapp-devstream-demo
    branch: main
    image_repo: ironcore864/go-webapp-devstream-demo
- name: githubactions-golang
  instanceID: default
  dependsOn: ["repo-scaffolding.golang-github"]
  options:
    owner: ${{repo-scaffolding.golang-github.outputs.owner}}
    repo: ${{repo-scaffolding.golang-github.outputs.repo}}
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
        repository: ${{repo-scaffolding.golang-github.outputs.repo}}
```

In the example above:

- We put `repo-scaffolding.golang-github` as dependency by using the `dependsOn` keyword.
- We used `repo-scaffolding.golang-github`'s output as input for the `githubactions-golang` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.var}}` is the syntax for using an output.
