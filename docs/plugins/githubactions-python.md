# githubactions-python Plugin

This plugin creates Python GitHub Actions workflows.

## Usage

_This plugin depends on an environment variable "GITHUB_TOKEN". Set it before using this plugin._

If you don't know how to create this token, check out:
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "githubactions-python.yaml"
```

All parameters are mandatory.
