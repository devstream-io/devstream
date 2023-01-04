# jenkins-pipeline Plugin

This plugin is used to create a `Jenkins Pipeline` for Github/Gitlab Repo.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

``` yaml
--8<-- "jenkins-pipeline.yaml"
```

**Notes:**

- `scm` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).
- The `pipeline` config option controls `CI` stages; you can refer to [Pipeline Config](./pipeline.md) for more info.
- `jenkins.token` is the password of `jenkins.user`.