## 1 `githubactions-python` Plugin

This plugin creates Python GitHub Actions workflows.

## 2 Usage:

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

If you don't know how to create this token, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

```yaml
tools:
- name: python-demo-app
  plugin:
    # name of the plugin
    kind: githubactions-python
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
      name: python
      # version of the language
      version: "3.8"
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

All parameters are mandatory.
