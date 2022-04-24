package plugin

var GithubActionsGolangDefaultConfig = `tools:
# name of the tool
- name: githubactions-golang
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below to an existing repo.
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME    
    repo: YOUR_REPO_NAME
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
        enable: False
        # go test profile subcommand, OPTIONAL, the given value below is default value
        profile: "-race -covermode=atomic"
        output: "coverage.out"
    docker:
      # docker build/push related, default to false
      enable: True
      registry:
        # dockerhub or harbor, default to dockerhub
        type: dockerhub
        # dockerhub/harbor username
        username: YOUR_DOCKERHUB_USERNAME
        # dockerhub/harbor image repository name
        repository: YOUR_DOCKERHUB_REPOSITORY`
