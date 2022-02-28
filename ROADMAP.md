# Roadmap

---

## 1 User-Facing Features Roadmap

### New Tools to Support (Plugins)

- Repository bootstrapping for Python/Nodejs for GitHub
- Repository bootstrapping for Golang/Python/Nodejs for GitLab
- GitLab CI workflows for Python/Nodejs
- Jenkins-GitHub, Jenkins-GitLab integration
- Jira-GitHub, Jira-GitLab integration
- ArgoCD-GitHub SSO integration
- FluxCD plugin
- Trello plugin that creates boards
- Trello-GitLab Integration

### Add a "Destroy" or "Force Delete" Feature

- Everything can be cleared up without any residue or side effects.

### Add an "Output" Feature

- All plugin's output is printed for users to review.

---

## 2 Technical Roadmap

### Plugins Dependency Management

- A common way to handle plugin dependencies and execution order using graph/topology sort or other applicable algorithms.
- Parallel/concurrency for plugins.

### Automated End-to-End Testing in a Staging Environment

- AWS EKS cluster creation with Terraform
- AWS EC2 (linux-amd64) creation with Terraform/Ansible
- Push notification to Slack/Lark when the testing environment is occupied or released

### Plugin Storage

Background: currently, we use GitHub releases to store pre-built binaries and plugins. We might need to figure out a better way to store/distribute plugins and binaries.

- Consider AWS S3 or similar choices for plugin storage.
- Make sure people who don't have optimum internet connections (e.g., users behind firewall or proxy) can still use DevStream smoothly.

### Output

- Consider adding "output" for each plugin, so that the dependant can use the output of the dependency (related to the user-facing roadmap, too.)

### Misc

- Shorter CI time: for example, adding packages into a base image.
- Send push notification to core committers when there is a pull request ready for review.

### Automated Release

- When a new tag is generated, build binaries for different platforms/OS and distribute the binaries to the plugin storage.

---

## 3 Already Done
