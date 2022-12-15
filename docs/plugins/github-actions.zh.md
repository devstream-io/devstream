# github-actions 插件

_该插件用于在项目中创建 Github Action Workflows_

## 用例

_该插件运行依赖以下环境变量：_

- GITHUB_TOKEN

请在使用插件前配置该环境变量。如果你不知道如何创建这个 TOKEN，可以查看以下文档：

- [创建个人访问 token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

下面的内容是一个示例配置文件用于创建 Github Workflows：

``` yaml
--8<-- "github-actions.yaml"
```

该插件的 `pipeline` 选项具体配置可查询[pipline配置项](pipeline.zh.md)。