# The State Section in the Main Config

In the main config, we can specify which "backend" to use to store DevStream state. Example:

Main config file:

```yaml
varFile: variables.yaml

toolFile: tools.yaml

state:
  backend: local
  options:
    stateFile: devstream.state
```

The `state` section has a `backend` parameter. At the moment, the only supported backend is the `local` one. AWS S3 will be supported soon.

For the `local` backend type, there is one option that is `stateFile`, which is the `path/to/the/state/file.state`.
