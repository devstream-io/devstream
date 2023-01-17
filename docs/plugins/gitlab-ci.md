# gitlab-ci Plugin

This plugin creates GitLab CI workflow.

It downloads a template of your choice, renders it with provided parameters, and creates a GitLab CI file for your repo.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file, and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

Plugin config example:

```yaml
--8<-- "gitlab-ci.yaml"
```

**Notes:**

- `scm` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- The `pipeline` config option controls `CI` stages; you can refer to [Pipeline Config](./pipeline.md) for more info.
