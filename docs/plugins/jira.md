# jira Plugin

This plugin integrates Jira with your GitHub repo.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "jira.yaml"
```

**Notes:**

- Jira language must be English
- There should be an existing Jira project
- `scm` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- `jira.token` should be created before using this plugin; you can refer to [Manage API tokens for your Atlassian account](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/).
