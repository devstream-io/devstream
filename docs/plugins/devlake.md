# `devlake` Plugin

This plugin installs [DevLake](https://github.com/merico-dev/lake) in an existing K8s cluster.

Note that this isn't a production-ready installation; it's only meant as an alternative for [the original docker-compose installation](https://github.com/merico-dev/lake/blob/main/docker-compose.yml) should you prefer K8s.

## Usage

```yaml
tools:
- name: devlake
  # name of the plugin
  plugin: devlake
  # if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
```

All the parameters in the example above are mandatory.
