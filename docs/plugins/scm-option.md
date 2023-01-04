# SCM Config Option

SCM Config Option is used to represent codebase-related config option.

## Config Options

| Option                | Description              |                    Note                                      |
| --------------------- | -----------------        | ------------------------------------------                   |
| owner                 | the owner of repo        | This option can't be configured at the same time with `org`  |
| org                   | the organization of repo | This option can't be configured at the same time with `owner`|
| scmType               | the repo type            | Support `Gitlab`/`Github` for now                            |
| name                  | the repo name            |                                                              |
| baseURL               | the gitlab url address   | If you use `Gitlab` for SCM, you should set this field       |
| url                   | the repo url address     | If you configure this option, then `org`, `owner`, `scmType`, `name` field can be empty |
| token                 | the repo api token       |                                                              |
| branch                | the repo branch          | If this option is empty, For `Github`, branch will be `main`, for `Gitlab`, branch will be `master` |

**Notes:**

_You need to get the token of the repo first._

- For `Github` Repo, you can refer to this [doc](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) about how to get `Github` token.
- For gitlab.com (instead of a self-hosted GitLab), [click here](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api) to create a token for DevStream (the scope contains API only).
- For self-hosted GitLab, refer to the [official doc here](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more info.

## Example

### Github SCM Config With URL

```yaml
scm:
  url: https://github.com/devstream-io/dtm-repo-scaffolding-python-flask.git
  branch: main
  token: TEST_TOKEN
```

### Github SCM Config Without URL

```yaml
scm:
  name: dtm-repo-scaffolding-python-flask
  scmType: github
  org: devstream-io
  branch: main
  token: TEST_TOKEN
```

### Gitlab SCM Config With URL

```yaml
scm:
  url: https://test.gitlab.com/testUser/dtm-repo-scaffolding-python-flask.git
  branch: master
  token: TEST_TOKEN
```

### Gitlab SCM Config With URL

```yaml
scm:
  name: dtm-repo-scaffolding-python-flask
  owner: testUser
  baseURL: https://test.gitlab.com
  branch: master
  token: TEST_TOKEN
  scmType: gitlab
```