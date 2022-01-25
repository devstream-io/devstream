## 1 GitHub Actions Plugin

This plugin creates some GitHub Actions workflows.

Currently, only Golang is supported.

We will support Python/Node.js soon, according to our roadmap.

## 2 Usage:

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
