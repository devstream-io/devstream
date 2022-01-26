## 1 GitHub Actions Plugin

This plugin creates some GitHub Actions workflows.

Currently, Golang, Python, and Node.js are supported.

## 2 Usage:

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

```yaml
tools:
# name of the plugin
- name: githubactions
  # version of the plugin
  version: 0.0.1
  # options for the plugin
  # checkout the version from the GitHub releases
  options:
    # the repo's owner
    owner: ironcore864
    # the repo where you'd like to setup GitHub Actions
    repo: go-hello-http
    # programming language specific settings
    language:
      # currently only go is supported
      name: go
      # version of the language
      version: "1.17"
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: master
    # configurations for the pipeline in GitHub Actions
    jobs:
      build:
        enable: True
        # build command, default to "go build ./..."
        command: "go build ./..."
      test:
        enable: True
        # test command, default to "go test ./..."
        command: "go test ./..."
        coverage:
          enable: True
          profile: "-race -covermode=atomic"
          output: "coverage.out"
      # docker build/push related
      docker:
        enable: True
        # dockerhub image repo
        repo: golang-demo
```

Currently, all the parameters in the example above are mandatory.
