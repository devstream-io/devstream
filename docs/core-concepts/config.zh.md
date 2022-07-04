# DevStream配置

DevStream使用YAML文件来描述你的DevOps工具链配置。

## 主配置文件

默认情况下，`dtm` 会使用当前目录下的`config.yaml`来作为主配置。

主配置包含三个部分：

- `varFile`: 指向定义变量的文件路径
- `toolFile`: 指向定义插件的文件路径
- `state`: 指定DevStream状态存储位置

### 主配置文件示例

`config.yaml` 的结构通常如下：

```yaml
varFile: variables.yaml

toolFile: tools.yaml

state:
  backend: local
  options:
    stateFile: devstream.state
```

### varFile

变量文件是一个用`key: value`格式来定义变量的YAML文件。

_At the moment, nested/composite values (for example, the value is a list/dictionary) are not supported yet._

`variables.yaml` 的结构通常如下:

```yaml
githubUsername: daniel-hutao
repoName: dtm-test-go
defaultBranch: main
dockerhubUsername: exploitht
```

### toolFile

插件文件是一个包含多种插件的Yaml文件。

- 插件文件只包含一个`tools`
- `tools`是一个定义多个插件的列表
- 列表中的每个对象都定义了一个由DevStream插件管理的工具
    - `name`: 是一个不带下划线的字符串，用来定义插件的名称
    - `instanceID`: 插件id
    - 你可以在一个插件文件中重复定义`name`，也可以在一个插件文件中重复定义`instanceID`，但是`name + instanceID`组合在一个插件文件中必须是唯一的
- 每个插件都有一个可选字段，即“选项”，它又是一个包含该特定插件参数的字典。关于插件的参数，请参见本文档的“插件”部分
- 每个插件都有一个可选字段，即“dependsOn”。继续阅读有关依赖项的详细信息。

`tools.yaml` 的结构通常如下:

```yaml
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: [[ githubUsername ]]
    org: ""
    repo: [[ repoName ]]
    branch: [[ defaultBranch ]]
    image_repo: [[ dockerhubUsername ]]/[[ repoName ]]
- name: jira-github-integ
  instanceID: default
  dependsOn: [ "github-repo-scaffolding-golang.default" ]
  options:
    owner: [[ githubUsername ]]
    repo: [[ repoName ]]
    jiraBaseUrl: https://xxx.atlassian.net
    jiraUserEmail: foo@bar.com
    jiraProjectKey: zzz
    branch: main
```

### state

`state`用来指定DevStream状态存储的位置，v0.5.0以前，DevStream仅支持状态记录存放在本地。

从v0.6.0开始，我们将支持`local`和`s3`两种存储。

更多状态存储细节请参见[DevStream状态存储](./stateconfig.md)

## 默认值

默认，`dtm` 使用 `config.yaml` 来作为主配置文件

### 指定主配置文件

你可以通过`dtm -f` or `dtm --config-file`来指定主配置文件。例如：

```shell
dtm apply -f path/to/your/config.yaml
dtm apply --config-file path/to/your/config.yaml
```

### varFile和toolFile默认没有值

对于`varFile`和`toolFile`, 默认没有任何值。

如果主配置中没有指定`varFile`，`dtm`将不会使用任何var文件，即使当前目录下已经有一个名为`variables.yaml`的文件。

同样，如果主配置中没有指定`toolFile`，即使当前目录下有`tools.yaml`文件，`dtm`也会抛出错误。
