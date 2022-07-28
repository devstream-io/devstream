# gitlabci-java plugin

This plugin set up Gitlab Pipeline in an existing Gitlab Java repository.

## Usage

**Notes:**

1. This plugin requires a Java Gitlab repository first

2. If `Build` is enabled (see the example below), you need to set `DOCKERHUB_TOKEN` environment variable. This will push the new built image to your repository(only support docker hub now).

3. If `Deploy` is enabled, you need to offer the Gitlab Kubernetes agent name(see [Gitlab-Kubernetes](https://docs.gitlab.cn/jh/user/clusters/agent/) for more details). This will deploy the new built application to your Kubernetes cluster. This step will use `deployment.yaml` to automatically deploy the application. Please create `manifests` directory in the repository root and create your `deployment.yaml` configuration file in it.

```yaml

--8<-- "gitlabci-java.yaml"

```
