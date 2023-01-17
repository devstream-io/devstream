# 用 DevStream 替代 Helm 让应用部署更加简单

helm-installer 插件实现了比 helm 更加简单和容易上手的方式来快速部署提供了 Helm Chart 的应用。

下面以 Argo CD 为例，介绍如何使用 DevStream 部署 Argo CD。

## 1 下载

进入你的工作目录，运行：

```shell
sh -c "$(curl -fsSL https://download.devstream.io/download.sh)"
```

!!! quote "可选"
    你可以将 `dtm` 移到 PATH 中。例如：`mv dtm /usr/local/bin/`。

    更多安装方式详见[安装dtm](../install.zh.md)。

## 2 配置

创建一个 `config.yaml` 文件，内容如下：

```yaml title="config.yaml"
config:
  state:
    backend: local
    options:
      stateFile: devstream.state
tools:
- name: helm-installer
  instanceID: argocd
```

其中 `instanceID` 为 "argocd"，匹配了 "argocd" 前缀，DevStream 会识别这个前缀，尝试寻找 Argo CD 应用对应的 Chart，并设置一系列默认值，然后开始部署。

通过 DevStream 安装 Helm 应用，你不需要搜索/阅读应用的官方文档，也不需要依次运行 `helm repo add` 等命令。你只需要知道应用的名称，将其作为 `instanceID` 的前缀，然后运行即可。这里是 [DevStream 支持的所有应用列表及前缀对应关系](../plugins/helm-installer/helm-installer.zh.md#3)。

## 3 初始化

运行以下命令以下载相应的插件：

```shell
./dtm init -f config.yaml
```

## 4 应用

运行以下命令以安装 Argo CD：

```shell
./dtm apply -f config.yaml -y
```

<script id="asciicast-549701" src="https://asciinema.org/a/549701.js" async></script>

## 5 检查结果

运行以下命令，可以看到 Argo CD 已经安装成功：

```shell
kubectl get pods -n argocd
```

<script id="asciicast-549703" src="https://asciinema.org/a/549703.js" async></script>

## 6 更进一步

DevStream 除了通过提供默认配置来简化应用部署，还提供了完整的 Helm values 文件的配置能力，详见[自定义Chart配置](../plugins/helm-installer/helm-installer.zh.md#22-chart)
