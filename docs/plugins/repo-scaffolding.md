# repo-scaffolding Plugin

This plugin bootstraps a GitHub or GitLab repo with scaffolding code for a web application.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file, and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "repo-scaffolding.yaml"
```

**Notes:**

- If you run `dtm delete`, the repo will be completely removed.
- If the `Update` interface is called, the repo will be completely removed and recreated.
- For the `repo-scaffolding` plugin, we only need `repo`, `delete_repo` permission for the token.
- `destinationRepo` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- `sourceRepo` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- if `destinationRepo` type is `Gitlab`, then `Devstream` supports config `destinationRepo`.`visibility`. This configuration is used to set the permissions of the new repository. The options are `public`, `private`, and `internal`.


### vars

This configuration is used for template rendering, It has default variables listed below:

```json
{
    "AppName": destinationRepo.repo,
    "Repo": {
        "Name": destinationRepo.repo,
        "Owner": destinationRepo.owner
    }
}
```

## Outputs

This plugin has three outputs:

- `owner`
- `repo`
- `repoURL`


## Examples 

### official scaffolding repo config

These repositories are official scaffolding repo to use for `sourceRepo` config; You can use these repo directly or just create one for yourself.

| language    | org           | repo                                |
|-------------|---------------|-------------------------------------|
| Golang      | devstream-io  | dtm-repo-scaffolding-golang-gin     |
| Golang      | devstream-io  | dtm-repo-scaffolding-golang-cli     |
| Python      | devstream-io  | dtm-repo-scaffolding-python-flask   |
| Java        | devstream-io  | dtm-repo-scaffolding-java-springboot|


### Golang

```yaml
tools:
- name: repo-scaffolding
  instanceID: golang-scaffolding
  options:
    destinationRepo:
      owner: test_owner
      org: ""
      name: dtm-test-golang
      branch: main
      scmType: github
    sourceRepo:
      org: devstream-io
      name: dtm-repo-scaffolding-golang-gin
      scmType: github
    vars:
      imageRepo: dtm-test/golang-repo
```

This config will create `dtm-test-golang` repo for user test_owner in GitHub, and the variable ImageRepo will be used for template rendering. 

### Golang CLI

```yaml
tools:
- name: repo-scaffolding
  instanceID: golang-cli-scaffolding
  options:
    destinationRepo:
      owner: test_owner
      org: ""
      name: dtm-test-golang-cli
      branch: main
      scmType: github
    sourceRepo:
      org: devstream-io
      name: dtm-repo-scaffolding-golang-cli
      scmType: github
```

This config will create `dtm-test-golang-cli` repo for user test_owner in GitHub.

### Java Spring

```yaml
tools:
- name: repo-scaffolding
  instanceID: java-scaffolding
  options:
    destinationRepo:
      owner: test_owner
      org: ""
      name: dtm-test-java
      branch: main
      baseUrl: 127.0.0.1:30001
      visibility: public
      scmType: gitlab
    sourceRepo:
      org: devstream-io
      name: dtm-repo-scaffolding-java-springboot
      scmType: github
```

This config will create `dtm-test-java` repo for user test_owner in GitHub.
