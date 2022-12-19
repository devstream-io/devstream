# GitOps

## 0 目标
在本教程中，我们将尝试通过 DevStream 来实现以下目标：

1. 创建一个 Python Web 应用程序仓库，基于 [Flask](https://flask.palletsprojects.com/en/2.2.x/) 框架；
2. 使用 GitHub Actions 为我们创建的仓库设置基本的 CI 流水线；
3. 在 _一个已有的 Kubernetes 集群_ 中安装 [Argo CD](https://argo-cd.readthedocs.io/en/stable/) 以实现 GitOps；
4. 创建一个 Argo CD 应用程序，用于部署第 1 步中生成的 Web 应用程序。

> 注意：
> 
> 在第 3 步中，Argo CD 安装在一个已有的 Kubernetes 集群中。DevStream 不配置基础设施，例如 Kubernetes 集群。
> 
> 如果你想跟着本教程自己尝试一下，但不知道如何在本地启动和运行 Kubernetes 集群，下面的博客（也来自 DevStream）可能会有所帮助：
> 
> - [用 Kind 从零开始快速搭建本地 Kubernetes 测试环节](https://blog.devstream.io/posts/%E7%94%A8kind%E9%83%A8%E7%BD%B2k8s%E7%8E%AF%E5%A2%83/)
> - [minikube结合阿里云镜像搭建本地开发测试环境](https://blog.devstream.io/posts/%E4%BD%BF%E7%94%A8minikube%E5%92%8C%E9%98%BF%E9%87%8C%E4%BA%91%E9%95%9C%E5%83%8F%E5%AE%89%E8%A3%85k8s/)

---

## 1 太长不看版：Demo 演示

如果你想看看 GitOps 的实际运行效果，可以看看下面的视频演示：

<iframe width="100%" height="500" src="https://www.youtube.com/embed/q7TK3vFr1kg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

这个演示录制于 DevStream 的旧版本，配置文件略有不同，但你还是能从中领略 DevStream GitOps 流程的魅力与要点。我们会尽快更新 DevStream 的最新版本的视频演示，敬请期待~

对于中文读者，可以看看这个：

<iframe src="//player.bilibili.com/player.html?aid=426762434&bvid=BV1W3411P7oW&cid=728576152&high_quality=1&danmaku=0" allowfullscreen="allowfullscreen" width="100%" height="500" scrolling="no" frameborder="0" sandbox="allow-top-navigation allow-same-origin allow-forms allow-scripts"></iframe>

光看完全不尽兴吧？跟着后面的步骤一起试一试吧！

---

## 2 概览

DevStream 将使用下面的插件来实现[第 0 节](#)中描述的目标：

1. [repo-scaffolding](../plugins/repo-scaffolding.md)
2. [github-actions](../plugins/github-actions.md)
3. [helm-installer](../plugins/helm-installer/helm-installer.md)
4. [argocdapp](../plugins/argocdapp.md)

不过，你不需要担心这些插件，因为 DevStream 会帮你自动管理它们。

---

## 3 启程：下载 DevStream (`dtm`)

为本教程创建一个临时工作目录：

```bash
mkdir test
cd test/
```

接着，在新创建的目录下，运行下面的命令：

```shell
sh -c "$(curl -fsSL https://download.devstream.io/download.sh)
```

这个脚本会根据你的操作系统来下载对应的 `dtm` 二进制文件。然后，赋予其可执行权限。

如果你执行 `ls` 命令，你会看到 `dtm` 二进制文件已经被下载下来了：

```bash
tiexin@mbp ~/work/devstream-io/test $ ls
dtm
```

然后，出于测试目的，我们可以尝试运行它，你会看到类似下面的输出：

```bash
tiexin@mbp ~/work/devstream-io/test $ ./dtm
DevStream is an open-source DevOps toolchain manager

######                 #####
#     # ###### #    # #     # ##### #####  ######   ##   #    #
#     # #      #    # #         #   #    # #       #  #  ##  ##
#     # #####  #    #  #####    #   #    # #####  #    # # ## #
#     # #      #    #       #   #   #####  #      ###### #    #
#     # #       #  #  #     #   #   #   #  #      #    # #    #
######  ######   ##    #####    #   #    # ###### #    # #    #

Usage:
  dtm [command]

Available Commands:
  apply       Create or update DevOps tools according to DevStream configuration file
  completion  Generate the autocompletion script for dtm for the specified shell
  delete      Delete DevOps tools according to DevStream configuration file
  destroy     Destroy DevOps tools deployment according to DevStream configuration file & state file
  develop     Develop is used for develop a new plugin
  help        Help about any command
  init        Download needed plugins according to the config file
  list        This command only supports listing plugins now
  show        Show is used to print plugins' configuration templates or status
  upgrade     Upgrade dtm to the latest release version
  verify      Verify DevOps tools according to DevStream config file and state
  version     Print the version number of DevStream

Flags:
      --debug   debug level log
  -h, --help    help for dtm

Use "dtm [command] --help" for more information about a command.
```

> 可选：你可以把 `dtm` 移动到 $PATH 环境变量中的某个目录下。例如：`mv dtm /usr/local/bin/`。这样，你就可以直接运行 `dtm` 而不需要再加上 `./` 前缀了。
> 
> 更多安装方式详见[安装 dtm](../install.zh.md)。

---

## 4 配置文件

运行以下命令来生成 gitops 的模板配置文件 `config.yaml` 。

```shell
./dtm show config -t gitops > config.yaml
```

按需修改 `config.yaml` 文件中的 `vars` 部分。记得修改 `githubUser` 和 `dockerUser` 的值为你自己的用户名。

在上面的例子中，我把这些变量设置成了下面的值：

| 变量        | 例子         | 说明                                |
|------------|-------------|------------------------------------|
| githubUser | IronCore864 | 大小写敏感，请改成你的 GitHub 用户名    |
| dockerUser | ironcore864 | 大小写敏感，请改成你的 DockerHub 用户名 |

## 5 环境变量

我们还需要设置以下环境变量：

```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
export IMAGE_REPO_PASSWORD="YOUR_DOCKERHUB_TOKEN_HERE"
```

> 提示：
> 如果你不知道如何创建这两个 token，可以参考：
> 
> - GITHUB_TOKEN：[Manage API tokens for your Atlassian account](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
> - IMAGE_REPO_PASSWORD：[Manage access tokens](https://docs.docker.com/docker-hub/access-tokens/)

---

## 6 初始化（Init）

运行以下命令，以根据配置文件自动下载所需插件：

```bash
./dtm init -f config.yaml
```

你会看到类似下面的输出：

```bash
2022-12-05 17:46:01 ℹ [INFO]  Using dir </Users/tiexin/.devstream/plugins> to store plugins.
2022-12-05 17:46:01 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.1  ] --------------------
... (略)
... (略)
2022-12-05 17:46:51 ✔ [SUCCESS]  Initialize finished.
```

---

## 7 应用（Apply）

运行：

```bash
./dtm apply -f config.yaml -y
```

你会看到类似下面的输出：

```
2022-12-05 17:49:49 ℹ [INFO]  Apply started.
2022-12-05 17:49:49 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-05 17:49:49 ℹ [INFO]  Tool (repo-scaffolding/myapp) found in config but doesn't exist in the state, will be created.
2022-12-05 17:49:49 ℹ [INFO]  Tool (helm-installer/argocd) found in config but doesn't exist in the state, will be created.
2022-12-05 17:49:49 ℹ [INFO]  Tool (github-actions/flask) found in config but doesn't exist in the state, will be created.
2022-12-05 17:49:49 ℹ [INFO]  Tool (argocdapp/default) found in config but doesn't exist in the state, will be created.
2022-12-05 17:49:49 ℹ [INFO]  Start executing the plan.
2022-12-05 17:49:49 ℹ [INFO]  Changes count: 4.
... (略)
... (略)
2022-12-05 17:51:51 ℹ [INFO]  -------------------- [  Processing progress: 4/4.  ] --------------------
2022-12-05 17:51:51 ℹ [INFO]  Processing: (argocdapp/default) -> Create ...
2022-12-05 17:51:52 ℹ [INFO]  application.argoproj.io/helloworld created
2022-12-05 17:51:52 ✔ [SUCCESS]  Tool (argocdapp/default) Create done.
2022-12-05 17:51:52 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-05 17:51:52 ✔ [SUCCESS]  All plugins applied successfully.
2022-12-05 17:51:52 ✔ [SUCCESS]  Apply finished.
```

---

## 8 查看结果

让我们来看看 `apply` 命令的结果。

### 8.1 GitHub 仓库

DevStream 已经通过 `repo-scaffolding` 插件自动创建了一个仓库：

![](gitops/a.png)

### 8.2 基于 GitHub Actions 的 CI 流水线

GitHub Actions 流水线已经被创建并运行：

![](gitops/b.png)

### 8.3 Argo CD 的安装

Argo CD 已经被安装到了 Kubernetes 集群中：

```bash
tiexin@mbp ~/work/devstream-io/test $ kubectl get namespaces
NAME                 STATUS   AGE
argocd               Active   5m42s
default              Active   6m28s
kube-node-lease      Active   6m29s
kube-public          Active   6m29s
kube-system          Active   6m29s
local-path-storage   Active   6m25s
tiexin@mbp ~/work/devstream-io/test $ kubectl get pods -n argocd
NAME                                               READY   STATUS    RESTARTS   AGE
argocd-application-controller-0                    1/1     Running   0          5m43s
argocd-applicationset-controller-66687659f-dsrtd   1/1     Running   0          5m43s
argocd-dex-server-6944757486-clshl                 1/1     Running   0          5m43s
argocd-notifications-controller-7944945879-b9878   1/1     Running   0          5m43s
argocd-redis-7887bbdbbb-xzppj                      1/1     Running   0          5m43s
argocd-repo-server-d4f5cc7cb-8gj24                 1/1     Running   0          5m43s
argocd-server-5bb75c4bd9-g948r                     1/1     Running   0          5m43s
```

### 8.4 使用 Argo CD 持续部署

CI 流水线已经构建了一个 Docker 镜像并推送到了 Dockerhub，而 DevStream 创建的 Argo CD 应用也部署了这个应用：

```bash
tiexin@mbp ~/work/devstream-io/test $ kubectl get deployment -n default
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
helloworld   1/1     1            1           5m16s
tiexin@mbp ~/work/devstream-io/test $ kubectl get pods -n default
NAME                          READY   STATUS    RESTARTS   AGE
helloworld-69b5586b94-wjwd9   1/1     Running   0          5m18s
tiexin@mbp ~/work/devstream-io/test $ kubectl get services -n default
NAME         TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)   AGE
helloworld   ClusterIP   10.96.73.97   <none>        80/TCP    5m27s
kubernetes   ClusterIP   10.96.0.1     <none>        443/TCP   8m2s
```

我们可以通过端口转发来访问这个应用：

```bash
kubectl port-forward -n default svc/helloworld 8080:80
```

在浏览器中访问 `localhost:8080`，你可以看到应用返回了一个 "Hello, World!"。大功告成！

---

## 9 清理

运行：

```bash
./dtm delete -f config.yaml -y
```

你会看到如下的输出：

```bash
2022-12-05 17:59:25 ℹ [INFO]  Delete started.
2022-12-05 17:59:26 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-05 17:59:26 ℹ [INFO]  Tool (argocdapp/default) will be deleted.
2022-12-05 17:59:26 ℹ [INFO]  Tool (github-actions/flask) will be deleted.
2022-12-05 17:59:26 ℹ [INFO]  Tool (repo-scaffolding/myapp) will be deleted.
2022-12-05 17:59:26 ℹ [INFO]  Tool (helm-installer/argocd) will be deleted.
2022-12-05 17:59:26 ℹ [INFO]  Start executing the plan.
2022-12-05 17:59:26 ℹ [INFO]  Changes count: 4.
... (略)
... (略)
2022-12-05 17:59:35 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-05 17:59:35 ✔ [SUCCESS]  All plugins deleted successfully.
2022-12-05 17:59:35 ✔ [SUCCESS]  Delete finished.
```

后面我们就能删除创建的所有文件了：

```bash
cd ../
rm -rf test/
rm -rf ~/.devstream/
```
