# The State Section in the Main Config

In the main config, we can specify which "backend" to use to store DevStream state.

We support the following types of backend:

- local
- s3

## Local backend config example:

```yaml
varFile: variables-gitops.yaml

toolFile: tools-gitops.yaml

state:
  backend: local
  options:
    stateFile: devstream.state
```

The `stateFile` under the `options` section is mandatory for local backend.

## S3 backend config example:

```yaml
varFile: variables-gitops.yaml

toolFile: tools-gitops.yaml

state:
  backend: s3
  options:
    bucket: devstream-remote-state
    region: ap-southeast-1
    key: devstream.state
```

The `bucket`, `region`, and `key` under the `ptions` section are all mandatory fields for s3 backend.
