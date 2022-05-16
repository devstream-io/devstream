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
tools:
# name of the tool
- name: gitlab-repo-scaffolding-golang
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitLab user name; please change the value below.
    owner: YOUR_GITLAB_USERNAME
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    # the repo which you'd like to create; please change the value below.
    repo: YOUR_REPO_NAME
    # the branch of the repo you'd like to hold the code
    branch: main
    # the image repo you'd like to push the container image; please change the value below.
    image_repo: YOUR_DOCKERHUB_USERNAME/YOUR_DOCKERHUB_REPOSITORY
```

Replace the following from the config above:

- `YOUR_GITLAB_USERNAME`
- `YOUR_ORGANIZATION_NAME`
- `YOUR_REPO_NAME`
- `YOUR_DOCKERHUB_USERNAME`
- `YOUR_DOCKERHUB_REPOSITORY`

The "branch" in the example above is "main", but you can adjust accordingly.

You have to specify either "owner" or "org".

## Outputs

This plugin has three outputs:

- `org`
- `owner`
- `repo`
- `repoURL` (example: "https://gitlab.com/ironcore864/dtm-app-test.git")