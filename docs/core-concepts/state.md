# State

## 1 Concept

State records the current status of the DevOps platform defined, created and managed by DevStream. DevStream relies on the State (and config, for that matter) to calculate required actions to ensure your DevOps platform is the same as defined in the config.

A `backend` is where to store the state, which we can configure in the config. At the moment, the following types of backends are supported:

- local
- k8s
- s3

## 2 How to Config the State

In the `config.state` section of the config, we can define where and how to store DevStream state.

### 2.1 Local File

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state # optional, defaults to "devstream.state"
```

### 2.2 Kubernetes

```yaml
config:
  state:
    backend: k8s
    options:
      namespace: devstream # optional, defaults to "devstream", create if not exist
      configmap: state     # optional, defaults to "state", create if not exist
```

### 2.3 S3

```yaml
config:
  state:
    backend: s3
    options:
      bucket: devstream-remote-state
      region: ap-southeast-1
      key: devstream.state
```

_Note: `options` `bucket`、`region` and `key` under the options are mandatory keys for s3 backend._
