# repo-scaffolding Plugin

This plugin bootstraps a GitHub or GitLab repo with scaffolding code for a web application.

## Requirement

This plugin need fllowing config base on your repo type:

### GitHub

- GITHUB_TOKEN: Set it before using this plugin.If you don't know how to create this token, check out [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

### GitLab

- GITLAB_TOKEN: Please set the environment variable before using the plugin. If you do not know how to create the token, Can view the document [Personal access tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).

- `destinationRepo.baseUrl`: If you are using a self-built GitLab repository, set this configuration to the URL of the self-built GitLab.

- `destinationRepo.visibility`: This configuration is used to set the permissions of the new repository. The options are `public`, `private`, and `internal`.

*Tips:*

- If you run `dtm delete`, the repo will be completely removed.

- If the `Update` interface is called, the repo will be completely removed and recreated. 

- For the  `repo-scaffolding` plugin, we only need `repo`, `delete_repo` permission for the token.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "repo-scaffolding.yaml"
```

Replace the following from the config above:

### destinationRepo

This configuration is used for the target repo, it includes the following config.

- `YOUR_DESTINATION_USERNAME`
- `YOUR_DESTINATION_ORGANIZATION_NAME`
- `YOUR_DESTINATION_REPO_NAME`
- `YOUR_DESTINATION_REPO_MAIN_BRANCH`
- `YOUR_DESTINATION_REPO_TYPE` 

**Please note that all parameter is case-sensitive.**

Currently, `owner`, `org`, and `repo` are mandatory, `branch` has the default value "main", `repoType` support  `gitlab` and `github` for now. 

### sourceRepo

This configuration is used for the source scaffolding repoI(only supports GitHub for now). It includes the following config.

- `YOUR_TEMPLATE_REPO_ORG`
- `YOUR_TEMPLATE_REPO_NAME`
- `YOUR_TEMPLATE_REPO_TYPE`

All the parameters in the example above are mandatory for now, `repoType` only support `github` for now. 

### vars

This configuration is used for template render, It has default variables listed below:

```json
{
    "AppName": destinationRepo.repo,
    "Repo": {
        "Name": destinationRepo.repo,
        "Owner": destinationRepo.owner
    }
}
```

## Examples 

### official scaffolding repo config

These repos are official scaffolding repo to use for `sourceRepo` config, You can use these repo directly or just create one for yourself.

| language | org | repo |
|  ----  | ----  |----  |
| Golang | devstream-io | dtm-scaffolding-golang |
| Java Spring | spring-guides | gs-spring-boot |


### Golang

```yaml
tools:
  - name: repo-scaffolding
    instanceID: golang-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        repo: dtm-test-golang
        branch: main
        repoType: github
      sourceRepo:
        org: devstream-io
        repo: dtm-scaffolding-golang
        repoType: github
      vars:
        ImageRepo: dtm-test/golang-repo
```

This config will create `dtm-test-golang` repo for user test_owner in GitHub, and the variable ImageRepo will be used for template render. 

### Java Spring

```yaml
tools:
  - name: repo-scaffolding
    instanceID: java-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        repo: dtm-test-java
        branch: main
        baseUrl: 127.0.0.1:30001
        visibility: public
        repoType: gitlab
      sourceRepo:
        org: spring-guides
        repo: gs-spring-boot
        repoType: github
```

this config will create `dtm-test-java` repo for user test_owner in GitHub.

## Outputs

This plugin has three outputs:

- `owner`
- `repo`
- `repoURL`
