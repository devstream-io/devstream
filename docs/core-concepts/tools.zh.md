# 工具(Tools)

## 1 概览

DevStream 视一切为 "工具"：

- 每个 工具 对应一个 DevStream 插件，它可以安装、配置或集成一些 DevOps 工具。
- 每个 工具 都有它的名称（对应插件名称）、实例 ID 和选项(Options)。
- 每个 工具 都可以依赖其他的工具，通过 `dependsOn` 关键字指定。

依赖 `dependsOn` 是一个字符串数组，每个元素都是一个依赖项，格式是 "工具名.实例ID"。

## 2 配置方式

DevStream 通过在配置中定义 `tools` 来声明所需的工具集合： 

- `tools` 是一个定义多个 工具 的列表
- 列表中的每个对象都定义了一个由 DevStream 插件管理的工具
    - `name`: 是一个不带下划线的字符串，用来定义插件的名称
    - `instanceID`: 实例的 ID，唯一标识一个工具实例
    - `name` 和 `instanceID` 可以分别重复，但是 `name + instanceID` 的组合必须唯一 
- 每个插件都有一个可选字段，"选项"(`options`)，每个插件的选项都是不同的，详情请参考[插件列表](../plugins/plugins-list.md)。
- 每个插件都有一个可选字段，"依赖项"(`dependsOn`)，定义了该插件依赖的其他插件列表。如果 A 依赖了 B 和 C，那么 dtm 会在 B 和 C 插件执行成功后再执行 A 插件。

`tools` 的配置示例如下：

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

其中，[[ githubUsername ]]、[[ repoName ]] 等是全局变量，它们的值可以在 `vars` 字段中定义。

