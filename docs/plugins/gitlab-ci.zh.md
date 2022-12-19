# gitlab-ci 插件

_该插件用于在使用 gitlab 的项目仓库中创建 Gitlab CI 并运行对应项目的 gitlab runner_

## 用例

_该插件运行依赖以下环境变量：_

- GITLAB_TOKEN

请在使用插件前配置该环境变量。如果你不知道如何创建这个 TOKEN，可以查看以下文档：

- 如果你使用的是官方提供的仓库（非自搭建的 `Gitlab`），可以查看该[链接](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api) 来创建一个用于 Devstream 的 token。
- 如果你使用自搭建的 `Gitlab` 仓库，可以查看该[链接](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api)来创建一个用于 Devstream 的 token。

### 配置项
下面的内容是一个示例配置文件用于创建 Gitlab CI：

``` yaml
--8<-- "gitlab-ci.yaml"
```

该插件的 `pipeline` 选项具体配置可查询[pipline配置项](pipeline.zh.md)。
