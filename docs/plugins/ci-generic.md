# ci-generic Plugin

This plugin installs ci script in GitLib/GitHub repo from local or remote url.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

``` yaml
--8<-- "ci-generic.yaml"
```

### Field Configs

| key                    | description                                                                                                      |
| ----                   | ----                                                                                                             |
| ci.localPath           | If your ci file is local, you can set the this field to the ci file location, which can be a directory or a file |
| ci.remoteURL           | If your ci file is remote, you can set this field to url address                                                 |
| ci.type                | ci type, support gitlab, github, jenkins for now                                                                 |
| projectRepo.owner      | destination repo owner                                                                                           |
| projectRepo.org        | destination repo org                                                                                             |
| projectRepo.name       | destination repo name                                                                                            |
| projectRepo.branch     | destination repo branch                                                                                          |
| projectRepo.scmType    | destination repo type, support github/gitlab for now                                                             |
| projectRepo.baseURL    | if you use self-build gitlab, you can set this field to gitlab address                                           |
| projectRepo.visibility | if you use gitlab, you can set this field for repo permission                                                    |

**Notes:**

- `ci.localPath` and `ci.remoteURL` can't be empty at the same time.
- if you set both `ci.localPath` and `ci.remoteURL`, `ci.localPath` will be used.
- if your `projectRepo.scmType` is `gitlab`, the `ci.type` is not allowed to be `github`.
- if your `projectRepo.scmType` is `github`, the `ci.type` is not allowed to be `gitlab`.

### Example

#### Local WorkFlows With Github

```yaml
tools:
- name: ci-generic
  instanceID: test-github
  options:
    ci:
      localPath: workflows
      type: github
    projectRepo:
      owner: devstream
      org: ""
      name: test-repo
      branch: main
      scmType: github
```

This config will put local workflows directory to GitHub repo's .github/workflows directory.

#### Remote Jenkinsfile With Gitlab

```yaml
tools:
- name: ci-generic
  instanceID: test-gitlab
  options:
    ci:
      remoteURL : https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile
      type: jenkins
    projectRepo:
      owner: root
      org: ""
      name: test-repo
      branch: main
      scmType: gitlab
      baseURL: http://127.0.0.1:30000
```

This config will put file from [remote](https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile)  to GitLab repo.
