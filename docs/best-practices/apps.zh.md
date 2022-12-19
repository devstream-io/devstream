# 使用 DevStream 的新特性 Apps 来实现 GitOps

## 0 目标

在本教程中，我们会使用 DevStream 的新特性 应用（Apps），来达到与 [GitOps](gitops.zh.md) 相似的效果。但它的配置更短，来展示 应用 的强大能力。如果你还没有读过原始的 GitOps 最佳实践，可以先点击前面的链接。

我们会创建两个应用程序（一个基于 Python，另一个是 Go 语言），并且创建共用的 CI/CD 流水线，即两个应用程序都会通过 Argo CD 来部署，就像前面的 GitOps 做到的那样。

---

## 1 概览

DevStream 使用 应用 的概念来创建带有代码脚手架和 GitHub Actions 流水线的仓库，并安装 Argo CD，最后使用 Arco CD 部署应用程序。

---

## 2 创建配置文件

为本教程创建临时工作目录：

```bash
mkdir test
cd test/
```

下载 dtm（详见 [GitOps](./gitops.zh.md) 最佳实践，如果你还没有下载过的话）

运行以下命令以生成配置文件：

```bash
./dtm show config --template=apps > config.yaml
```

替换下面命令中的双引号内里面的内容，并运行，以设置环境变量：

```shell
export GITHUB_USER="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export DOCKERHUB_USERNAME="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

现在我们就可以使用前面设置的环境变量来更新配置文件了：

===  "基于 **macOS** 或 **FreeBSD** 的系统"

    ```shell
    sed -i.bak "s@YOUR_GITHUB_USER@${GITHUB_USER}@g" config.yaml
    sed -i.bak "s@YOUR_DOCKERHUB_USER@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

=== "**GNU** Linux 用户"

    ```shell
    sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
    sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

---

## 3 初始化（Init）和应用（Apply）

运行下面这个命令，以根据配置文件来自动下载对应的插件：

```bash
./dtm init -f config.yaml
```

再运行 apply 命令：

```bash
./dtm apply -f config.yaml -y
```

你会看到类似下面的输出：

```bash
tiexin@mbp ~/work/devstream-io/test $ ./dtm apply -f config.yaml -y
2022-12-16 14:30:27 ℹ [INFO]  Apply started.
2022-12-16 14:30:28 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-16 14:30:28 ℹ [INFO]  Tool (helm-installer/argocd) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (repo-scaffolding/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (repo-scaffolding/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (github-actions/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (argocdapp/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (github-actions/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (argocdapp/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Start executing the plan.
2022-12-16 14:30:28 ℹ [INFO]  Changes count: 7.
2022-12-16 14:30:28 ℹ [INFO]  -------------------- [  Processing progress: 1/7.  ] --------------------
2022-12-16 14:30:28 ℹ [INFO]  Processing: (helm-installer/argocd) -> Create ...
2022-12-16 14:30:29 ℹ [INFO]  Filling default config with instance: argocd.
2022-12-16 14:30:29 ℹ [INFO]  Creating or updating helm chart ...
... (略)
... (略)
2022-12-16 14:32:09 ℹ [INFO]  -------------------- [  Processing progress: 7/7.  ] --------------------
2022-12-16 14:32:09 ℹ [INFO]  Processing: (argocdapp/myapp2) -> Create ...
2022-12-16 14:32:19 ℹ [INFO]  application.argoproj.io/myapp2 created
2022-12-16 14:32:19 ✔ [SUCCESS]  Tool (argocdapp/myapp2) Create done.
2022-12-16 14:32:19 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-16 14:32:19 ✔ [SUCCESS]  All plugins applied successfully.
2022-12-16 14:32:19 ✔ [SUCCESS]  Apply finished.
```

---

## 4 查看结果

让我们来继续看看 `apply` 命令的结果：

与我们在[GitOps](./gitops.zh.md)最佳实践中所做的类似，我们可以检查 dtm 是否为两个应用程序创建了代码仓库并为其创建了 CI 流水线，同时安装了 Argo CD，而且两个应用程序都使用了 Argo CD 来部署到了 Kubernetes 集群中。

---

## 5 清理

运行:

```bash
./dtm delete -f config.yaml -y
```

你会看到类似下面的输出：

```bash
2022-12-16 14:34:30 ℹ [INFO]  Delete started.
2022-12-16 14:34:31 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-16 14:34:31 ℹ [INFO]  Tool (github-actions/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (argocdapp/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (github-actions/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (argocdapp/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (repo-scaffolding/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (repo-scaffolding/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (helm-installer/argocd) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Start executing the plan.
2022-12-16 14:34:31 ℹ [INFO]  Changes count: 7.
2022-12-16 14:34:31 ℹ [INFO]  -------------------- [  Processing progress: 1/7.  ] --------------------
2022-12-16 14:34:31 ℹ [INFO]  Processing: (github-actions/myapp1) -> Delete ...
2022-12-16 14:34:33 ℹ [INFO]  Prepare to delete 'github-actions_myapp1' from States.
2022-12-16 14:34:33 ✔ [SUCCESS]  Tool (github-actions/myapp1) delete done.
... (略)
... (略)
2022-12-16 14:34:40 ✔ [SUCCESS]  Tool (helm-installer/argocd) delete done.
2022-12-16 14:34:40 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-16 14:34:40 ✔ [SUCCESS]  All plugins deleted successfully.
2022-12-16 14:34:40 ✔ [SUCCESS]  Delete finished.
```

后面我们就能删除创建的所有文件了：

```bash
cd ../
rm -rf test/
rm -rf ~/.devstream/
```
