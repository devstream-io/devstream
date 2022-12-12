# 快速开始

我们将在本文使用 DevStream 自动完成以下操作：

- 创建一个包含了 web 应用程序的 GitHub 仓库，代码基于 [gin](https://github.com/gin-gonic/gin) 框架（用Go语言编写）自动生成；
- 为前面创建的仓库设置 GitHub Actions 工作流。

---

## 1 下载

进入你的工作目录，运行：

```shell
sh -c "$(curl -fsSL https://download.devstream.io/download.sh)"
```

!!! note "提示"
    上面的命令会做以下事情：

    - 检测你的操作系统和芯片架构
    - 找到最新版本的 `dtm` 二进制文件
    - 根据操作系统和架构下载正确的 `dtm` 二进制文件
    - 授予二进制文件执行权限

!!! quote "可选"
    你可以将 `dtm` 移到 PATH 中。例如：`mv dtm /usr/local/bin/`。

    更多安装方式详见[安装dtm](./install.zh.md)。

---

## 2 配置

运行以下命令来生成 quickstart 的模板配置文件 `config.yaml` 。

```shell
./dtm show config -t quickstart > config.yaml
```

运行以下命令以设置这些环境变量（记得替换双引号内的值）：

```shell
export GITHUB_USER="<YOUR_GITHUB_USER_NAME_HERE>"
export GITHUB_TOKEN="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export DOCKERHUB_USERNAME="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

!!! tip "提示"
    前往 [Personal Access Token](https://github.com/settings/tokens/new) 为 `dtm` 生成新的 `GITHUB_TOKEN`。

    对于“快速开始”，我们只需要勾选 `repo`、`workflow`、`delete_repo` 权限。

接着，让我们运行以下命令，以使用环境变量来修改配置文件：

===  "**macOS** 或 基于 **FreeBSD** 的操作系统"

    ```shell title=""
    sed -i.bak "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
    sed -i.bak "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

=== "**GNU** Linux 用户"

     ```shell title=""
     sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
     sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
     ```

---

## 3 初始化

运行：

```shell
./dtm init -f config.yaml
```

!!! success "你会看到类似下面的输出"
    ``` title=""
    2022-12-12 11:40:34 ℹ [INFO]  Using dir </Users/stein/.devstream/plugins> to store plugins.
    2022-12-12 11:40:34 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:35 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.2.so] ...
     87.75 MiB / 87.75 MiB [================================] 100.00% 7.16 MiB/s 12s
    2022-12-12 11:40:47 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.2.so] download succeeded.
    2022-12-12 11:40:48 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.2.md5] ...
     33 B / 33 B [=========================================] 100.00% 115.84 KiB/s 0s
    2022-12-12 11:40:48 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.2.md5] download succeeded.
    2022-12-12 11:40:48 ℹ [INFO]  Initialize [repo-scaffolding-darwin-arm64_0.10.2] finished.
    2022-12-12 11:40:48 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.2  ] --------------------

    2022-12-12 11:40:48 ℹ [INFO]  -------------------- [  github-actions-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:48 ℹ [INFO]  Downloading: [github-actions-darwin-arm64_0.10.2.so] ...
     90.27 MiB / 90.27 MiB [================================] 100.00% 10.88 MiB/s 8s
    2022-12-12 11:40:57 ✔ [SUCCESS]  [github-actions-darwin-arm64_0.10.2.so] download succeeded.
    2022-12-12 11:40:57 ℹ [INFO]  Downloading: [github-actions-darwin-arm64_0.10.2.md5] ...
     33 B / 33 B [=========================================] 100.00% 145.46 KiB/s 0s
    2022-12-12 11:40:57 ✔ [SUCCESS]  [github-actions-darwin-arm64_0.10.2.md5] download succeeded.
    2022-12-12 11:40:57 ℹ [INFO]  Initialize [github-actions-darwin-arm64_0.10.2] finished.
    2022-12-12 11:40:57 ℹ [INFO]  -------------------- [  github-actions-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:57 ✔ [SUCCESS]  Initialize finished.
    ```

---

## 4 应用（Apply）

运行：

```shell
./dtm apply -f config.yaml -y
```

!!! success "你会看到类似下面的输出"

    ```text title=""
    2022-12-12 11:44:39 ℹ [INFO]  Apply started.
    2022-12-12 11:44:39 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-12 11:44:39 ℹ [INFO]  Tool (repo-scaffolding/golang-github) found in config but doesn't exist in the state, will be created.
    2022-12-12 11:44:39 ℹ [INFO]  Tool (github-actions/default) found in config but doesn't exist in the state, will be created.
    2022-12-12 11:44:39 ℹ [INFO]  Start executing the plan.
    2022-12-12 11:44:39 ℹ [INFO]  Changes count: 2.
    2022-12-12 11:44:39 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-12 11:44:39 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Create ...
    2022-12-12 11:44:39 ℹ [INFO]  github start to download repoTemplate...2022-12-12 11:44:42 ✔ [SUCCESS]  The repo go-webapp-devstream-demo has been created.
    2022-12-12 11:44:49 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) Create done.
    2022-12-12 11:44:49 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-12 11:44:49 ℹ [INFO]  Processing: (github-actions/default) -> Create ...
    2022-12-12 11:44:57 ✔ [SUCCESS]  Tool (github-actions/default) Create done.
    2022-12-12 11:44:57 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-12 11:44:57 ✔ [SUCCESS]  All plugins applied successfully.
    2022-12-12 11:44:57 ✔ [SUCCESS]  Apply finished.
    ```

---

## 5 检查结果

前往你的 GitHub 仓库列表，可以看到一个新的仓库 `go-webapp-devstream-demo` 已经被创建了。

包含了 Golang web 应用程序的脚手架代码，并正确设置了 GitHub Actions CI 工作流。

DevStream 在生成仓库脚手架和创建工作流时的代码提交，已经触发了 CI，且工作流已经成功地运行完毕，如下图所示：

![](./images/quickstart.png)

---

## 6 清理

运行：

```shell
./dtm delete -f config.yaml
```

输入 `y` 然后回车，你会看到类似下面的输出：

!!! success "输出"
    ```title=""
    2022-12-12 12:29:00 ℹ [INFO]  Delete started.
    2022-12-12 12:29:00 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-12 12:29:00 ℹ [INFO]  Tool (github-actions/default) will be deleted.
    2022-12-12 12:29:00 ℹ [INFO]  Tool (repo-scaffolding/golang-github) will be deleted.
    Continue? [y/n]
    Enter a value (Default is n): y
    2022-12-12 12:29:00 ℹ [INFO]  Start executing the plan.
    2022-12-12 12:29:00 ℹ [INFO]  Changes count: 2.
    2022-12-12 12:29:00 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-12 12:29:00 ℹ [INFO]  Processing: (github-actions/default) -> Delete ...
    2022-12-12 12:29:02 ℹ [INFO]  Prepare to delete 'github-actions_default' from States.
    2022-12-12 12:29:02 ✔ [SUCCESS]  Tool (github-actions/default) delete done.
    2022-12-12 12:29:02 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-12 12:29:02 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Delete ...
    2022-12-12 12:29:03 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
    2022-12-12 12:29:03 ℹ [INFO]  Prepare to delete 'repo-scaffolding_golang-github' from States.
    2022-12-12 12:29:03 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) delete done.
    2022-12-12 12:29:03 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-12 12:29:03 ✔ [SUCCESS]  All plugins deleted successfully.
    2022-12-12 12:29:03 ✔ [SUCCESS]  Delete finished.
    ```

现在，如果你看看 GitHub 仓库列表，所有东西都被 DevStream 消灭了。妙哉！

你也可以通过运行 `rm devstream.state` 来删除 DevStream 状态文件（现在应该是个空文件）。
