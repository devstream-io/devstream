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

In the following tool file example, tool "github-repo-scaffolding-golang" (with instance id "default") will be installed before tool "githubactions-golang" (with instance id "default"):

```yaml
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    org: devstream-io
    repo: dtm-e2e-go
    branch: main
    image_repo: dtme2etest/dtm-e2e-go
- name: githubactions-golang
  instanceID: default
  dependsOn: ["github-repo-scaffolding-golang.default"]
  options:
    org: ${{github-repo-scaffolding-golang.default.outputs.org}}
    repo: ${{github-repo-scaffolding-golang.default.outputs.repo}}
    language:
      name: go
      version: "1.17"
    branch: main
    build:
      enable: True
    test:
      enable: True
      coverage:
        enable: True
```
