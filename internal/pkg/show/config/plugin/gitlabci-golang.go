package plugin

var GitlabCIGolangDefaultConfig = `tools:
- name: go-hello-world
  # name of the plugin
  plugin: gitlabci-golang
  # options for the plugin
  options:
    # owner/repo; "path with namespace" is only GitLab API's way of saying the same thing; please change the values below.
    pathWithNamespace: YOUR_GITLAB_USERNAME/YOUR_GITLAB_REPO_NAME
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main`
