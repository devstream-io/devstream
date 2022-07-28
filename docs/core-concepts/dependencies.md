## Tool Dependencies

If you want tool A to be installed before tool B, you can let tool B depend on tool A.

The syntax for dependency is:
    
```yaml
dependsOn: [ "ToolName.ToolInstanceID" ]
```

Since `dependsOn` is a list, a tool can have multiple dependencies:

```
dependsOn: [ "ToolName1.ToolInstanceID1", "ToolName2.ToolInstanceID2", "..." ]
```

In the following tool file example, tool "repo-scaffolding" (with instance id "golang-github") will be installed before tool "githubactions-golang" (with instance id "default"):

```yaml
tools:
- name: repo-scaffolding
  instanceID: golang-github
  options:
    destination_repo:
      owner: [[ githubUsername ]]
      org: ""
      repo: [[ repoName ]]
      branch: [[ defaultBranch ]]
      repo_type: github
    vars:
      ImageRepo: "[[ dockerhubUsername ]]/[[ repoName ]]"
    source_repo:
      org: devstream-io
      repo: dtm-scaffolding-golang
      repo_type: github
- name: githubactions-golang
  instanceID: default
  dependsOn: ["repo-scaffolding.golang-github"]
  options:
    org: ${{repo-scaffolding.golang-github.outputs.org}}
    repo: ${{repo-scaffolding.golang-github.outputs.repo}}
    language:
      name: go
      version: "1.18"
    branch: main
    build:
      enable: True
    test:
      enable: True
      coverage:
        enable: True
```
