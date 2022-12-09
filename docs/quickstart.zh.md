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
    2022-12-02 16:11:55 ℹ [INFO]  Using dir </Users/tiexin/.devstream/plugins> to store plugins.
    2022-12-02 16:11:55 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.1  ] --------------------
    2022-12-02 16:11:57 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.1.so] ...
     87.82 MiB / 87.82 MiB [================================] 100.00% 12.30 MiB/s 7s
    2022-12-02 16:12:04 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.1.so] download succeeded.
    2022-12-02 16:12:04 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.1.md5] ...
     33 B / 33 B [==========================================] 100.00% 50.98 KiB/s 0s
    2022-12-02 16:12:04 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.1.md5] download succeeded.
    2022-12-02 16:12:04 ℹ [INFO]  Initialize [repo-scaffolding-darwin-arm64_0.10.1] finished.
    2022-12-02 16:12:04 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.1  ] --------------------
    2022-12-02 16:12:04 ℹ [INFO]  -------------------- [  githubactions-golang-darwin-arm64_0.10.1  ] --------------------
    2022-12-02 16:12:05 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.10.1.so] ...
     86.44 MiB / 86.44 MiB [================================] 100.00% 15.12 MiB/s 5s
    2022-12-02 16:12:10 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.10.1.so] download succeeded.
    2022-12-02 16:12:10 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.10.1.md5] ...
     33 B / 33 B [==========================================] 100.00% 71.24 KiB/s 0s
    2022-12-02 16:12:10 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.10.1.md5] download succeeded.
    2022-12-02 16:12:11 ℹ [INFO]  Initialize [githubactions-golang-darwin-arm64_0.10.1] finished.
    2022-12-02 16:12:11 ℹ [INFO]  -------------------- [  githubactions-golang-darwin-arm64_0.10.1  ] --------------------
    2022-12-02 16:12:11 ✔ [SUCCESS]  Initialize finished.
    ```

---

## 4 应用（Apply）

运行：

```shell
./dtm apply -f config.yaml -y
```

!!! success "你会看到类似下面的输出"

    ```text title=""
    2022-12-02 16:18:00 ℹ [INFO]  Apply started.
    2022-12-02 16:18:00 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-02 16:18:00 ℹ [INFO]  Tool (repo-scaffolding/golang-github) found in config but doesn't exist in the state, will be created.
    2022-12-02 16:18:00 ℹ [INFO]  Tool (githubactions-golang/default) found in config but doesn't exist in the state, will be created.
    2022-12-02 16:18:00 ℹ [INFO]  Start executing the plan.
    2022-12-02 16:18:00 ℹ [INFO]  Changes count: 2.
    2022-12-02 16:18:00 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-02 16:18:00 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Create ...
    2022-12-02 16:18:00 ℹ [INFO]  github start to download repoTemplate...
    2022-12-02 16:18:04 ✔ [SUCCESS]  The repo go-webapp-devstream-demo has been created.
    2022-12-02 16:18:12 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) Create done.
    2022-12-02 16:18:12 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-02 16:18:12 ℹ [INFO]  Processing: (githubactions-golang/default) -> Create ...
    2022-12-02 16:18:13 ℹ [INFO]  Creating GitHub Actions workflow pr-builder.yml ...
    2022-12-02 16:18:14 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml created.
    2022-12-02 16:18:14 ℹ [INFO]  Creating GitHub Actions workflow main-builder.yml ...
    2022-12-02 16:18:15 ✔ [SUCCESS]  Github Actions workflow main-builder.yml created.
    2022-12-02 16:18:15 ✔ [SUCCESS]  Tool (githubactions-golang/default) Create done.
    2022-12-02 16:18:15 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-02 16:18:15 ✔ [SUCCESS]  All plugins applied successfully.
    2022-12-02 16:18:15 ✔ [SUCCESS]  Apply finished.
    ```

---

## 5 检查结果

前往你的 GitHub 仓库列表，可以看到一个新的仓库 `go-webapp-devstream-demo` 已经被创建了。

包含了 Golang web 应用程序的脚手架代码，并正确设置了 GitHub Actions CI 工作流。

DevStream 在生成仓库脚手架和创建工作流时的代码提交，已经触发了 CI，且工作流已经成功地运行完毕，如下图所示：

![](./images/repo-scaffolding.png)

---

## 6 清理

运行：

```shell
./dtm delete -f config.yaml
```

输入 `y` 然后回车，你会看到类似下面的输出：

!!! success "输出"
    ```title=""
    2022-12-02 16:19:07 ℹ [INFO]  Delete started.
    2022-12-02 16:19:07 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-02 16:19:07 ℹ [INFO]  Tool (githubactions-golang/default) will be deleted.
    2022-12-02 16:19:07 ℹ [INFO]  Tool (repo-scaffolding/golang-github) will be deleted.
    Continue? [y/n]
    Enter a value (Default is n): y
    
    2022-12-02 16:19:08 ℹ [INFO]  Start executing the plan.
    2022-12-02 16:19:08 ℹ [INFO]  Changes count: 2.
    2022-12-02 16:19:08 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-02 16:19:08 ℹ [INFO]  Processing: (githubactions-golang/default) -> Delete ...
    2022-12-02 16:19:09 ℹ [INFO]  Deleting GitHub Actions workflow pr-builder.yml ...
    2022-12-02 16:19:09 ✔ [SUCCESS]  GitHub Actions workflow pr-builder.yml removed.
    2022-12-02 16:19:10 ℹ [INFO]  Deleting GitHub Actions workflow main-builder.yml ...
    2022-12-02 16:19:10 ✔ [SUCCESS]  GitHub Actions workflow main-builder.yml removed.
    2022-12-02 16:19:10 ℹ [INFO]  Prepare to delete 'githubactions-golang_default' from States.
    2022-12-02 16:19:10 ✔ [SUCCESS]  Tool (githubactions-golang/default) delete done.
    2022-12-02 16:19:10 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-02 16:19:10 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Delete ...
    2022-12-02 16:19:11 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
    2022-12-02 16:19:11 ℹ [INFO]  Prepare to delete 'repo-scaffolding_golang-github' from States.
    2022-12-02 16:19:11 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) delete done.
    2022-12-02 16:19:11 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-02 16:19:11 ✔ [SUCCESS]  All plugins deleted successfully.
    2022-12-02 16:19:11 ✔ [SUCCESS]  Delete finished.
    ```

现在，如果你看看 GitHub 仓库列表，所有东西都被 DevStream 消灭了。妙哉！

你也可以通过运行 `rm devstream.state` 来删除 DevStream 状态文件（现在应该是个空文件）。
