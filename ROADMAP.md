# Roadmap

## 1 New Tools to Support (Plugins)

- artifactory: https://github.com/devstream-io/devstream/issues/607
- Jenkins pipeline plugin: https://github.com/devstream-io/devstream/issues/582
- Cloudflare IP List Monitor: https://github.com/devstream-io/devstream/issues/560
- Use Validator for Config Validation: https://github.com/devstream-io/devstream/issues/558
- GitLab CE: https://github.com/devstream-io/devstream/issues/509
- zendao: https://github.com/devstream-io/devstream/issues/508
- Tekton
- Jira-GitHub, Jira-GitLab integration
- Jenkins-GitHub, Jenkins-GitLab integration
- ArgoCD-GitHub SSO integration
- Repository bootstrapping for Python/Nodejs for GitHub, Golang/Python/Nodejs for GitLab
- GitLab CI workflows for Python/Nodejs
- FluxCD plugin
- Trello-GitLab Integration

## 2 Core Features

General:

- single config profile support: https://github.com/devstream-io/devstream/issues/596
- make sure people who don't have optimum internet connections (e.g., users behind firewall or proxy) can still use DevStream smoothly.

### `dtm show config`

This is already supported, but we will improve the features of it, for example:

- show the default configs of multiple plugins that are used together
- interactive: user select plugin then show the default config

## 3 Quality of Life Improvements for Developers

### Automated End-to-End Testing in a Staging Environment

- AWS EC2 (linux-amd64) creation with Terraform/Ansible
- Push notification to Slack/Lark when the testing environment is occupied or released

### Misc

- Integrate the golangci-lint command in the makefile https://github.com/devstream-io/devstream/issues/632
- Push notification when CI failed: https://github.com/devstream-io/devstream/issues/636
- Shorter CI time: for example, adding packages into a base image
- More end-to-end tests coverage, to test more typical usecases and plugins
- Push notification to core committers when there is a new PR ready for review

## 4 Already Done

Core:
- local state support: https://github.com/devstream-io/devstream/issues/16
- pluginmanager module: https://github.com/devstream-io/devstream/issues/17
- statemanager supports concurrent map: https://github.com/devstream-io/devstream/issues/71
- pluginmanager download status bar: https://github.com/devstream-io/devstream/issues/88, https://github.com/devstream-io/devstream/issues/98
- multi-plugin instance support: https://github.com/devstream-io/devstream/issues/136
- force delete feature: https://github.com/devstream-io/devstream/issues/177, https://github.com/devstream-io/devstream/issues/278
- `verify` command: https://github.com/devstream-io/devstream/issues/252, https://github.com/devstream-io/devstream/issues/253
- output support: https://github.com/devstream-io/devstream/issues/324
- Homebrew support: https://github.com/devstream-io/devstream/issues/351, https://github.com/devstream-io/devstream/issues/372
- remote state: https://github.com/devstream-io/devstream/issues/378, https://github.com/devstream-io/devstream/issues/485
- autocomplete: https://github.com/devstream-io/devstream/issues/380
- generate default config: https://github.com/devstream-io/devstream/issues/383
- list all plugins: https://github.com/devstream-io/devstream/issues/384
- config supports global variables: https://github.com/devstream-io/devstream/issues/393
- plugin status: https://github.com/devstream-io/devstream/issues/401
- parallel download: https://github.com/devstream-io/devstream/issues/579
- plugins released to AWS S3 instead of GitHub releases page

Plugins:
- Jenkins plugin: https://github.com/devstream-io/devstream/issues/11
- GitLab CI plugin: https://github.com/devstream-io/devstream/issues/12
- GitHub Actions plugin for Nodejs and Python: https://github.com/devstream-io/devstream/issues/14
- GitHub Actions plugin supports test coverage: https://github.com/devstream-io/devstream/issues/133
- GitHub repo scaffolding for Golang plugin: https://github.com/devstream-io/devstream/issues/191, https://github.com/devstream-io/devstream/issues/520
- Prometheus/grafana plugin: https://github.com/devstream-io/devstream/issues/231
- Helm type plugins support values: https://github.com/devstream-io/devstream/issues/272
- openldap plugin: https://github.com/devstream-io/devstream/issues/284
- trello plugin: https://github.com/devstream-io/devstream/issues/307, https://github.com/devstream-io/devstream/issues/314
- GitLab CI generic plugin: https://github.com/devstream-io/devstream/issues/377
- Helm generic plugin: https://github.com/devstream-io/devstream/issues/424

Develop:
- cross-platform build: https://github.com/devstream-io/devstream/issues/21, https://github.com/devstream-io/devstream/issues/170
- end to end test: https://github.com/devstream-io/devstream/issues/50, https://github.com/devstream-io/devstream/issues/118
- logging level: https://github.com/devstream-io/devstream/issues/176
- parallel build: https://github.com/devstream-io/devstream/issues/361
- automated release: https://github.com/devstream-io/devstream/issues/364
- command to help contributors to generate scaffolding code: https://github.com/devstream-io/devstream/issues/454, https://github.com/devstream-io/devstream/issues/443, https://github.com/devstream-io/devstream/issues/436
- params validation improvement: https://github.com/devstream-io/devstream/issues/511
- editor config: https://github.com/devstream-io/devstream/issues/629

By versions:

v0.3.1:
- automated release: when a new tag is generated, build binaries for different platforms/OS and distribute the binaries to the plugin storage.

v0.3.0:
- "Destroy" and "force delete": everything can be cleared up without any residue or side effects.
- "Output": all plugin's output is printed for users to review.
- Plugin dependency management: a common way to handle plugin dependencies and execution order using graph/topology sort.
- Automated e2e testing: AWS EKS cluster with Terraform.
- Trello plugin that creates boards.

v0.4.0-v0.5.0:
- generic GitLab CI plugin
- define variables and use it in the config file.
- auto-complete support for `dtm` commands
- HashiCorp Vault
