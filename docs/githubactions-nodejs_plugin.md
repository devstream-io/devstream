## 1 GitHub Actions Nodejs Plugin

This plugin creates Nodejs GitHub Actions workflows.

## 2 Usage:

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

```yaml
tools:
- name: nodejs-demo-app
  plugin:
    # name of the plugin
    kind: githubactions-nodejs
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.0.2
  # options for the plugin
  options:
    # the repo's owner
    owner: ironcore864
    # the repo where you'd like to setup GitHub Actions
    repo: nodejs-demo
    # programming language specific settings
    language:
      name: nodejs
      # version of the language
      version: "16.14"
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

All parameters are mandatory.
