# `devlake` Plugin

This plugin installs [DevLake](https://github.com/merico-dev/lake) in an existing K8s cluster.

Note that this isn't a production-ready installation; it's only meant as an alternative for [the original docker-compose installation](https://github.com/merico-dev/lake/blob/main/docker-compose.yml) should you prefer K8s.

## Usage

```yaml
tools:
- name: devlake
  # name of the plugin
  plugin: devlake
```

All the parameters in the example above are mandatory.
