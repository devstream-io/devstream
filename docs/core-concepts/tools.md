# Tools

## 1 Tools

DevStream treats everything as a concept named _Tool_:

- Each _Tool_ corresponds to a DevStream plugin, which can either install, configure, or integrate some DevOps tools.
- Each _Tool_ has its Name, InstanceID, and Options.
- Each _Tool_ can have its dependencies, specified by the `dependsOn` keyword.

The dependency `dependsOn` is an array of strings, each element being a dependency.

Each dependency is named in the format of "TOOL_NAME.INSTANCE_ID".

---

## 2 Configuration

Define your needed `tools` in DevStream config:

- `tools` is a list of `tool`.
- Each element in the list defines a DevOps tool (managed by a DevStream plugin), with the following key/values:.
    - `name`: a string without underscore, corresponds to the name of the plugin.
    - `instanceID`: unique instance ID of a tool.
    - Multiple tools defined with the same `name` or `instanceID` are allowd, but `name + instanceID` must be unique.
- Each plugin has an optional setting `options`, and the options for each plugin is different. See the [list of plugins](../plugins/plugins-list.md) for more details.
- Each plugin has an optional setting `dependsOn` which defines the dependencies of this plugin. E.g., if A depends on B and C, then dtm will only execute A after B and C.

An example of `tools` config:

```yaml
tools:
- name: repo-scaffolding
  instanceID: golang-github
  options:
    destinationRepo:
      owner: [[ githubUsername ]]
      name: [[ repoName ]]
      branch: [[ defaultBranch ]]
      scmType: github
    vars:
      ImageRepo: "[[ dockerhubUsername ]]/[[ repoName ]]"
    sourceRepo:
      org: devstream-io
      name: dtm-scaffolding-golang
      scmType: github
- name: jira-github-integ
  instanceID: default
  dependsOn: [ "repo-scaffolding.golang-github" ]
  options:
    owner: [[ githubUsername ]]
    repo: [[ repoName ]]
    jiraBaseUrl: https://xxx.atlassian.net
    jiraUserEmail: foo@bar.com
    jiraProjectKey: zzz
    branch: main
```

`[[ githubUsername ]]`, `[[ repoName ]]` (and other variables inside the double brackets) are global variables which are defined in the `vars` section of the config.
