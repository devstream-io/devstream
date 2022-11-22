# gitlabci-java 插件

## 用例

**注意:**

1. 使用该插件之前，需要在Gitlab上拥有一个Java代码仓库。

2. 如果`Build`选项被开启，需要设置`DOCKERHUB_TOKEN`环境变量。这会将新构建的镜像推送到你的镜像仓库（目前只支持docker hub仓库）。

3. 如果`Deploy`选项被开启，你需要提供Gitlab配置的Kubernetes代理名称（设置详情参照[Gitlab-Kubernetes](https://docs.gitlab.cn/jh/user/clusters/agent/))。这会将新构建的应用部署至上面提供的Kubernetes集群中。该步骤会使用`deployment.yaml`来自动部署应用，请在仓库根目录下创建`manifests`目录，并在其中新建你的`deployment.yaml`配置文件。

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/overview.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "gitlabci-java.yaml"
```
