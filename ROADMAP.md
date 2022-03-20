# Roadmap

## 1 User-Facing Features Roadmap

### New Tools to Support (Plugins)

- Jira-GitHub, Jira-GitLab integration
- Jenkins-GitHub, Jenkins-GitLab integration
- ArgoCD-GitHub SSO integration
- Repository bootstrapping for Python/Nodejs for GitHub
- Repository bootstrapping for Golang/Python/Nodejs for GitLab
- GitLab CI workflows for Python/Nodejs
- FluxCD plugin
- Trello-GitLab Integration

## 2 Technical Roadmap

### Automated End-to-End Testing in a Staging Environment

- AWS EC2 (linux-amd64) creation with Terraform/Ansible
- Push notification to Slack/Lark when the testing environment is occupied or released

### Automated Release

- When a new tag is generated, build binaries for different platforms/OS and distribute the binaries to the plugin storage.

### Plugin Storage

Background: currently, we use GitHub releases to store pre-built binaries and plugins. We might need to figure out a better way to store/distribute plugins and binaries.

- Consider AWS S3 or similar choices for plugin storage.
- Make sure people who don't have optimum internet connections (e.g., users behind firewall or proxy) can still use DevStream smoothly.

### Plugins Dependency Management

- Parallel/concurrency for plugins.

### Misc

- Shorter CI time: for example, adding packages into a base image.
- Send push notification to core committers when there is a pull request ready for review.

## 3 Already Done

v0.3.0:
- "Destroy" and "force delete": everything can be cleared up without any residue or side effects.
- "Output": all plugin's output is printed for users to review.
- Plugin dependency management: a common way to handle plugin dependencies and execution order using graph/topology sort.
- Automated e2e testing: AWS EKS cluster with Terraform.
- Trello plugin that creates boards.
