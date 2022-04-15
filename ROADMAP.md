# Roadmap

## 1 New Tools to Support (Plugins)

- Jira-GitHub, Jira-GitLab integration
- Jenkins-GitHub, Jenkins-GitLab integration
- HashiCorp Vault
- ArgoCD-GitHub SSO integration
- Repository bootstrapping for Python/Nodejs for GitHub, Golang/Python/Nodejs for GitLab
- generic GitLab CI plugin, GitLab CI workflows for Python/Nodejs
- FluxCD plugin
- Trello-GitLab Integration

## 2 `dtm` Core

### Variables

Define variables and use it in the config file.

### `dtm show state`

Check plugin(s)' state.

### `dtm show config`

This is already supported, but we will improve the features of it, for example:

- show the default config of one plugin
- show the default configs of multiple plugins that are used together
- interactive: user select plugin then show the default config

### Auto-Complete

Auto-complete support for `dtm` commands.

### Plugin Storage

Background: currently, we use GitHub releases to store pre-built binaries and plugins. We might need to figure out a better way to store/distribute plugins and binaries.

- Consider AWS S3 or similar choices for plugin storage.
- Make sure people who don't have optimum internet connections (e.g., users behind firewall or proxy) can still use DevStream smoothly.

## 3 Quality of Life Improvements

### Automated End-to-End Testing in a Staging Environment

- AWS EC2 (linux-amd64) creation with Terraform/Ansible
- Push notification to Slack/Lark when the testing environment is occupied or released

### Misc

- Shorter CI time: for example, adding packages into a base image
- More end-to-end tests coverage, to test more typical usecases and plugins
- Push notification to core committers when there is a new PR ready for review

## 4 Already Done

v0.3.1:
- automated release: when a new tag is generated, build binaries for different platforms/OS and distribute the binaries to the plugin storage.

v0.3.0:
- "Destroy" and "force delete": everything can be cleared up without any residue or side effects.
- "Output": all plugin's output is printed for users to review.
- Plugin dependency management: a common way to handle plugin dependencies and execution order using graph/topology sort.
- Automated e2e testing: AWS EKS cluster with Terraform.
- Trello plugin that creates boards.
