# ci-generic Plugin

This plugin installs `CI` script in `GitLib`/`GitHub` repo from a local or remote url.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

``` yaml
--8<-- "ci-generic.yaml"
```

**Notes:**

- `projectRepo` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- `ci.localPath` and `ci.remoteURL` can't be empty at the same time.
- if your `projectRepo.scmType` is `gitlab`, the `ci.type` is not allowed to be `github-actions`.
- if your `projectRepo.scmType` is `github`, the `ci.type` is not allowed to be `gitlab-ci`.

## Example

### Local WorkFlows With Github

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

### Remote Jenkinsfile With Gitlab

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

This config will put files from [remote](https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile)](https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile)  to GitLab repo.
