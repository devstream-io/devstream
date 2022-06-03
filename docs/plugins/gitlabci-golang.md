# gitlabci-golang Plugin

This plugin creates Golang GitLab CI workflow.

## Usage

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token._

TL;DR: if you are using gitlab.com (instead of a self-hosted GitLab), [click here](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api) to create a token for DevStream (the scope contains API only.)

If you are using self-hosted GitLab, refer to the [official doc here](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more info.

_Note: when creating the token, make sure you select "API" in the "scopes" section, as DevStream uses GitLab API to add CI workflow files._

Plugin config example:

```yaml
--8<-- "gitlabci-golang.yaml"
```

All parameters are mandatory.
