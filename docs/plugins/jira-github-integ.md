# jira-github-integ Plugin

This plugin integrates Jira with your GitHub repo.

## Usage

_Please confirm the preconditions:_

- Jira language must be English
- There should be an existing Jira project

_This plugin depends on the following two environment variables:_

- JIRA_API_TOKEN
- GITHUB_TOKEN

Set the values accordingly before using this plugin.

If you don't know how to create these tokens, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [Manage API tokens for your Atlassian account](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/).

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "jira-github-integ.yaml"
```

Currently, all the parameters in the example above are mandatory.
