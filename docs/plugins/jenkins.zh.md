# jenkins 插件

本插件使用 helm 在已有的 k8s 集群上安装 [Jenkins](https://jenkins.io)。

并且安装 [GitHub Pull Request Builder](https://plugins.jenkins.io/ghprb/) 插件和 [OWASP Markup Formatter](https://plugins.jenkins.io/antisamy-markup-formatter/) 插件；同时利用 OWASP Markup Formatter 插件激活 HTML 渲染模式。

## 使用方法

### 生产环境

请将配置中的 `storageClass` 修改为已存在的 StorageClass.

### 测试环境

如果你想**在本地测试插件**：

1. 请将配置文件中的 `test_env` 改为 `true`。
2. 在运行 Kubernetes 的主机上创建数据目录并修改权限，命令如下：

如果 Kubernetes 和 dtm 运行在同一个主机上：

```bash
mkdir -p ~/data/jenkins-volume/
chown -R 1000:1000 ~/data/jenkins-volume/
```

如果 Kubernetes 和 dtm 运行在不同的主机上，比如 Kubernetes 运行在 虚拟机或者 Docker 容器中：

```bash
# 1 获取 dtm 运行的主机的用户的 home 目录
cd ~ && pwd
# 2 复制上面的命令结果
# 3 进入运行着 k8s 的主机
  # 3.1 如果你是 minikube
  minikube ssh
  # 3.2 如果你是 kind
  docker exec -it <kind-container-name-or-id> bash
# 4. 在 k8s 运行的主机上创建数据目录并修改权限，命令如下：
mkdir -p <your-dtm-home-dir>/data/jenkins-volume/
chown -R 1000:1000 <your-dtm-home-dir>/data/jenkins-volume/
```

### 配置

```yaml
--8<-- "jenkins.yaml"
```

当前，所有配置项均为必填。

## 输出

这个插件有两个输出：

- `jenkinsURL` (格式: `hostname:port`, 例如: "localhost:8080")
- `jenkinsPasswordOfAdmin` 
