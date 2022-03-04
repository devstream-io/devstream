## 1 `githubactions-golang` Plugin

This plugin creates some Golang GitHub Actions workflows.

## 2 Usage:

_This plugin depends on the following environment variable:_

- GITHUB_TOKEN

Set it before using this plugin.

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
    version: 0.2.0
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
        enable: True
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
