# SCM 配置项

SCM 配置项用于表示代码仓库相关的配置信息。

## 配置项

| 字段                | 描述              |                    备注                                      |
| --------------------- | -----------------        | ------------------------------------------                   |
| owner                 | 仓库拥有者        | 这个字段不能和 `org` 同时配置  |
| org                   | 仓库所属的组织 | 这个字段不能和 `owner` 同时配置 |
| scmType               | 仓库类型            | 目前支持 `Gitlab`/`Github`                            |
| name                  | 仓库名           |                                                              |
| baseURL               | 仓库的根 url 地址   | 如果你使用的是 `Gitlab` 仓库，你就需要配置这个字段       |
| url                   | 仓库的 url 地址     | 如果你配置了这个字段, 那么 `org`, `owner`, `scmType`, `name` 就不需要配置 |
| token                 | 仓库的认证 token       |                                          |
| branch                | 仓库的分支         | 如果该字段为空，对于 `Github` 仓库，会使用 main 分支，对于 `Gitlab` 仓库，会使用 master 分支 |

**Notes:**

_你需要先获取仓库对应的 token。_

- 对于 `Github` 仓库，你可以查阅该[文档](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) 来获取 token。
- 对于 `Gitlab` 官方仓库（非自建）, 可以查看该[文档](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api)来创建 token (配置的 scope 只需要包含 API)。
- 对于自建的 `Gitlab` 仓库，可以查看该[文档](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more info。

## 示例

### 使用 url 配置 Github 仓库

```yaml
scm:
  url: https://github.com/devstream-io/dtm-repo-scaffolding-python-flask.git
  branch: main
  token: TEST_TOKEN
```

### 不使用 url 配置 Github 仓库

```yaml
scm:
  name: dtm-repo-scaffolding-python-flask
  scmType: github
  org: devstream-io
  branch: main
  token: TEST_TOKEN
```

### 使用 url 配置 Gitlab 仓库

```yaml
scm:
  url: https://test.gitlab.com/testUser/dtm-repo-scaffolding-python-flask.git
  branch: master
  token: TEST_TOKEN
```

### 不使用 url 配置 Gitlab 仓库

```yaml
scm:
  name: dtm-repo-scaffolding-python-flask
  owner: testUser
  baseURL: https://test.gitlab.com
  branch: master
  token: TEST_TOKEN
  scmType: gitlab
```