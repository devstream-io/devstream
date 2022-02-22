## 1 GitHub Actions Golang Plugin

This plugin creates some Golang GitHub Actions workflows.

## 2 Usage:

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

```yaml
tools:
- name: golang-demo-app
  plugin:
    # name of the plugin
    kind: githubactions-golang
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.0.2
  # options for the plugin
  options:
    # the repo's owner
    owner: ironcore864
    # the repo where you'd like to setup GitHub Actions
    repo: golang-demo
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
      # dockerhub image repo
      repo: golang-demo  
```

Some parameters are optional. See the default values and optional parameters in the example above.
