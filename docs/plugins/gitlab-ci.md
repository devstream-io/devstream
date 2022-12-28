# gitlab-ci Plugin

This plugin creates GitLab CI workflow.

It downloads a template of your choice, render it with provided parameters, and creates a GitLab CI file to your repo.

## Usage

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token._

TL;DR: if you are using gitlab.com (instead of a self-hosted GitLab), [click here](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api) to create a token for DevStream (the scope contains API only.)

If you are using self-hosted GitLab, refer to the [official doc here](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more info.

_Note: when creating the token, make sure you select "API" in the "scopes" section, as DevStream uses GitLab API to add CI workflow files._

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

Plugin config example:

```yaml
--8<-- "gitlab-ci.yaml"
```

Or, run `dtm show config --plugin=gitlab-ci` to get the default config.

All parameters are mandatory.
