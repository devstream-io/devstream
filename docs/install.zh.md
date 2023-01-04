# 安装

!!! note "提示"
    当前 dtm 发行版支持如下操作系统与 CPU 架构组合：

    1. Darwin/arm64
    2. Darwin/amd64
    3. Linux/amd64

## 1 用脚本安装

进入你的工作目录，运行：

```shell
sh -c "$(curl -fsSL https://download.devstream.io/download.sh)"
```

这个命令会根据你的操作系统和芯片架构下载对应的 `dtm` 二进制文件到你的工作目录中，并赋予二进制文件执行权限。

!!! note "可选"

    将 dtm 移动到包含于 PATH 的目录下，比如 `mv dtm /usr/local/bin/`。

## 2 从 GitHub Release 页面手动下载安装

在 [Release](https://github.com/devstream-io/devstream/releases/) 页面找到当前最新版本 `dtm`，然后点击下载。

如果你不方便使用浏览器下载，也可以选择使用 curl：

```shell
# 版本号 v0.10.3，系统类型 linux 和 CPU 架构 amd64 根据需要灵活修改
$ curl -o dtm https://download.devstream.io/v0.10.3/dtm-linux-amd64
```

需要注意的是当前 `dtm` 提供了多个版本，你需要根据操作系统和芯片架构选择自己需要的正确版本。下载到本地后，你可以选择将其重命名，移入包含在"$PATH"的目录里并赋予其可执行权限，比如在 Linux 上你可以执行如下命令完成这些操作：

```shell
mv dtm-linux-amd64 /usr/local/bin/dtm
chmod +x dtm
```

接着你可以通过如下命令验证 dtm 的权限以及版本等是否正确：

```shell
$ dtm version
0.10.3
```

## 3 用 [asdf](https://asdf-vm.com/) 安装

```shell
# 安装 dtm 插件
asdf plugin add dtm
# 列出所有可用版本
asdf list-all dtm
# 安装特定版本
asdf install dtm latest
# 设置全局版本 (在 ~/.tool-versions 文件中)
asdf global dtm latest
# 现在你就能使用 dtm 了
dtm help
```
