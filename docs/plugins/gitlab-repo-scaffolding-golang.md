# gitlab-repo-scaffolding-golang Plugin

This plugin bootstraps a GitLab repo with scaffolding code for a Golang web application.

_This plugin depends on the following environment variable:_

- GITLAB_TOKEN

Set it before using this plugin.

If you don't know how to create this token, check out:
- [GitLab Personal Access Tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

*Tips:*

*1. If you run `dtm delete`, the repo on GitLab will be completely removed.*

*2. If the `Update` interface is called, the repo on GitLab will be completely removed and recreated.

## Usage

**Please note that the `owner` parameter is case-sensitive.**

```yaml
--8<-- "gitlab-repo-scaffolding-golang.yaml"
```

Replace the following from the config above:

- `YOUR_GITLAB_USERNAME`
- `YOUR_ORG_OR_GROUP_NAME`
- `YOUR_REPO_NAME`
- `YOUR_DOCKERHUB_USERNAME`
- `YOUR_DOCKERHUB_REPOSITORY`

The "branch" in the example above is "main", but you can adjust accordingly.

You have to specify either "owner" or "org".

## Outputs

This plugin has four outputs:

- `org`
- `owner`
- `repo`
- `repoURL` (example: "https://gitlab.com/ironcore864/dtm-app-test.git")
