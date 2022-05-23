# 快速开始

> 中文 | [English](./quickstart_en.md)

## 1 下载 DevStream (`dtm`)

根据你的实际环境从 [DevStream Releases](https://github.com/devstream-io/devstream/releases) 下载合适的 `dtm` 版本。

> 记得将二进制文件重命名为 `dtm`，这样用起来更简单。比如：`mv dtm-darwin-arm64 dtm`。

> 一旦下载完成，你就可以将 dtm 文件放到任何目录下运行了。当然更加建议你将 dtm 加到 PATH 下(比如 `/usr/local/bin`)。

## 2 准备一个配置文件

开始之前：这是一个DevStream配置的例子：[examples/tools-quickstart.yaml](https://github.com/devstream-io/devstream/blob/main/examples/tools-quickstart.yaml) 。记得打开这个配置文件，把里面所有的 `FULL_UPPER_CASE_STRINGS`（比如说 `YOUR_GITHUB_USERNAME` ）修改成你自己的。注意每一项的含义，并确保它是你要的。对于其他插件，请查看我们的 [文档](https://docs.devstream.io) 中的"插件"部分，以了解详细用法。

将 [examples/quickstart.yaml](https://raw.githubusercontent.com/devstream-io/devstream/main/examples/quickstart.yaml) 和 [examples/tools-quickstart.yaml](https://raw.githubusercontent.com/devstream-io/devstream/main/examples/tools-quickstart.yaml) 文件下载到你到工作目录下，然后重命名`quickstart.yaml` 成 `config.yaml`：

```shell
curl -o config.yaml https://raw.githubusercontent.com/devstream-io/devstream/main/examples/quickstart.yaml
curl -o tools-quickstart.yaml https://raw.githubusercontent.com/devstream-io/devstream/main/examples/tools-quickstart.yaml
```

然后相应的修改配置文件中的内容。

比如我的 GitHub 用户名是 "IronCore864", 然后我的 Dockerhub 用户名是 "ironcore864"，这样我就可以运行：

```shell
sed -i.bak "s/YOUR_GITHUB_USERNAME_CASE_SENSITIVE/IronCore864/g" tools-quickstart.yaml

sed -i.bak "s/YOUR_DOCKER_USERNAME/ironcore864/g" tools-quickstart.yaml
```

> 这个配置文件会使用两个插件，一个用来创建 GitHub 项目，而且初始化成一个 Golang 的 web 应用结构。接着另外一个插件会给这个项目创建对应的 GitHub Actions 工作流。

这两个插件, [githubactions-golang](plugins/githubactions-golang.md)&[github-repo-scaffolding-golang](plugins/github-repo-scaffolding-golang.md), 需要配置一个环境变量 才能工作，我们看下怎么配置：


```shell
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
```

如果你不知道怎么创建一个 GitHub token 可以看下[官方文档](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) 。

## 3 初始化

运行：

```shell
dtm init -f config.yaml
```

然后你可以看到类似这样的日志输出：

```
2022-03-04 12:08:06 ℹ [INFO]  Initialize started.
2022-03-04 12:08:06 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-03-04 12:08:11 ℹ [INFO]  Downloading: [github-repo-scaffolding-golang-darwin-arm64_0.2.0.so] ...
 13.52 MiB / 13.52 MiB [=================================] 100.00% 3.14 MiB/s 4s
2022-03-04 12:08:15 ✔ [SUCCESS]  [github-repo-scaffolding-golang-darwin-arm64_0.2.0.so] download succeeded.
2022-03-04 12:08:17 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.2.0.so] ...
 16.05 MiB / 16.05 MiB [=================================] 100.00% 5.41 MiB/s 2s
2022-03-04 12:08:20 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.2.0.so] download succeeded.
2022-03-04 12:08:20 ✔ [SUCCESS]  Initialize finished.
```

此步骤验证您的 dtm 二进制文件的 MD5 sum，根据配置文件下载所需的插件，并验证插件的  MD5 sum。

注意：如果您的 dtm 二进制文件的 MD5 sum 与我们发布页面中的 MD5 sum 不匹配，dtm init 将停止。 如果您的本地 dtm MD5 不同，则表明您自己构建了二进制文件（例如，出于开发目的）。 由于 Go 插件的性质，dtm 必须与相应的插件一起构建。 所以，如果你正在构建 dtm，你也应该构建插件，在这种情况下，你不需要运行 dtm init 来下载插件。

## 4 命令用法总览

如果你需要应用配置，请运行：

```shell
./dtm apply -f YOUR_CONFIG_FILE.yaml
```

如果你没有用` -f `参数指定配置文件，它将尝试使用默认值，即当前目录下的 `config.yaml` 。

`dtm`将对比 `Config`、`State` 和 `Resource`，决定是否需要 `Create`、`Update` 或 `Delete`。更多信息请阅读我们的 [核心概念](https://docs.devstream.io/en/latest/core-concepts/core-concepts/) 文档。

上面的命令在实际执行改变之前会要求你确认。如果不需要确认就应用 `Config`（就像 `apt-get -y update` ），请运行：

```shell
./dtm -y apply -f YOUR_CONFIG_FILE.yaml
```

要删除 `Config` 中定义的所有内容，请运行:

```shell
./dtm delete -f YOUR_CONFIG_FILE.yaml
```

注意，这将删除 `Config` 中定义的所有内容。如果某些 `Config` 在应用后被删除（`State` 有，但 `Config` 没有），`dtm delete`不会删除它，这与`dtm destroy`不同。

同样的，如果不需要确认就删除内容，请运行：
```shell
./dtm -y delete -f YOUR_CONFIG_FILE.yaml
```

要删除 `Config` 中定义的所有内容，且无论 `State` 是什么：
```shell
./dtm delete --force -f YOUR_CONFIG_FILE.yaml
```

验证以上命令已正确执行，请运行：
```shell
./dtm verify -f YOUR_CONFIG_FILE.yaml
```

销毁所有内容，请运行：
```shell
./dtm destroy
```

`dtm`将读取 `State`，然后确定哪些 `Tool` 应被安装，然后删除这些 `Tool`。这与`dtm apply -f empty.yaml`相同（ `empty.yaml` 是一个空的配置文件）。

## 5 开始执行

运行：

```shell
dtm apply -f config.yaml
```

然后输入 y 来确认继续执行命令，接着你可以看到类似这样的日志输出：

```
2022-03-04 12:08:54 ℹ [INFO]  Apply started.
2022-03-04 12:08:54 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-03-04 12:08:54 ℹ [INFO]  Tool < go-webapp-repo > found in config but doesn't exist in the state, will be created.
2022-03-04 12:08:54 ℹ [INFO]  Tool < golang-demo-actions > found in config but doesn't exist in the state, will be created.
Continue? [y/n]
Enter a value (Default is n): y

2022-03-04 12:08:57 ℹ [INFO]  Start executing the plan.
2022-03-04 12:08:57 ℹ [INFO]  Changes count: 2.
2022-03-04 12:08:57 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-03-04 12:08:57 ℹ [INFO]  Processing: go-webapp-repo -> Create ...
2022-03-04 12:09:04 ℹ [INFO]  Repo created.
2022-03-04 12:09:22 ✔ [SUCCESS]  Plugin go-webapp-repo Create done.
2022-03-04 12:09:22 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-03-04 12:09:22 ℹ [INFO]  Processing: golang-demo-actions -> Create ...
2022-03-04 12:09:23 ℹ [INFO]  Language is: go-1.17.
2022-03-04 12:09:23 ℹ [INFO]  Creating GitHub Actions workflow pr-builder.yml ...
2022-03-04 12:09:24 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml created.
2022-03-04 12:09:25 ℹ [INFO]  Creating GitHub Actions workflow main-builder.yml ...
2022-03-04 12:09:26 ✔ [SUCCESS]  Github Actions workflow main-builder.yml created.
2022-03-04 12:09:26 ✔ [SUCCESS]  Plugin golang-demo-actions Create done.
2022-03-04 12:09:26 ✔ [SUCCESS]  All plugins applied successfully.
2022-03-04 12:09:26 ✔ [SUCCESS]  Apply finished.
```
## 6 检查结果

登录你自己的 GitHub 账户，然后你可以看到一个新的名字叫做 "go-webapp-devstream-demo" 的项目已经被创建出来了，
而且里面已经有了一些 Golang 的 web 应用脚手架代码。另外你还可以看到用于构建这个应用的一些 GitHub Actions 也已经被配置好了。酷吧？

## 7 清理

运行：

```shell
dtm destroy
```

然后你可以看到类似这样的日志：

```
2022-03-04 12:10:36 ℹ [INFO]  Destroy started.
2022-03-04 12:10:36 ℹ [INFO]  Change added: go-webapp-repo_github-repo-scaffolding-golang -> Delete
2022-03-04 12:10:36 ℹ [INFO]  Change added: golang-demo-actions_githubactions-golang -> Delete
Continue? [y/n]
Enter a value (Default is n): y

2022-03-04 12:10:38 ℹ [INFO]  Start executing the plan.
2022-03-04 12:10:38 ℹ [INFO]  Changes count: 2.
2022-03-04 12:10:38 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-03-04 12:10:38 ℹ [INFO]  Processing: go-webapp-repo -> Delete ...
2022-03-04 12:10:40 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
2022-03-04 12:10:40 ℹ [INFO]  Prepare to delete 'go-webapp-repo_github-repo-scaffolding-golang' from States.
2022-03-04 12:10:40 ✔ [SUCCESS]  Plugin go-webapp-repo delete done.
2022-03-04 12:10:40 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-03-04 12:10:40 ℹ [INFO]  Processing: golang-demo-actions -> Delete ...
2022-03-04 12:10:40 ℹ [INFO]  language is go-1.17.
2022-03-04 12:10:41 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml already removed.
2022-03-04 12:10:42 ✔ [SUCCESS]  Github Actions workflow main-builder.yml already removed.
2022-03-04 12:10:42 ℹ [INFO]  Prepare to delete 'golang-demo-actions_githubactions-golang' from States.
2022-03-04 12:10:42 ✔ [SUCCESS]  Plugin golang-demo-actions delete done.
2022-03-04 12:10:42 ✔ [SUCCESS]  All plugins destroyed successfully.
2022-03-04 12:10:42 ✔ [SUCCESS]  Destroy finished.
```
