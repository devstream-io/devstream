package config

var DefaultConfig = `# config.yaml sample:
# var file path, you can set it to absolute path or relative path.
varFile: variables.yaml # here is a relative path. (defaults is ./variables.yaml)
# tool file path, you can set it to absolute path or relative path.
toolFile: tools.yaml # here is a relative path.
# state config
state:
  backend: local # backend can be local or s3
  options:
    stateFile: devstream.state

# tools.yaml sample:
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: [[ githubUsername ]]
    org: ""
    repo: [[ repoName ]]

# variables.yaml sample:
githubUsername: daniel-hutao
repo: go-webapp-demo`
