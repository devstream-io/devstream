# 安装

## 0 当前支持的操作系统与芯片架构

*  Darwin/arm64
*  Darwin/amd64
*  Linux/amd64

## 1 用脚本安装

进入你的工作目录，运行：

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/install/download.sh)"
```

这个命令会根据你的操作系统和芯片架构下载对应的 `dtm` 二进制文件到你的工作目录中，并赋予二进制文件执行权限。

> 可选：建议你将 dtm 移动到包含于 PATH 的目录下，比如 `mv dtm /usr/local/bin/`。

## 2 用 [asdf](https://asdf-vm.com/) 安装

```shell
# Plugin
asdf plugin add dtm
# Show all installable versions
asdf list-all dtm
# Install specific version
asdf install dtm latest
# Set a version globally (on your ~/.tool-versions file)
asdf global dtm latest
# Now dtm commands are available
dtm --help
```

## 3 从 Release 页面手动下载

在 [Release](https://github.com/devstream-io/devstream/releases/) 页面找到当前最新版本 `dtm`，然后点击下载。
需要注意的是当前 `dtm` 提供了多个版本，你需要根据操作系统和芯片架构选择自己需要的正确版本。下载到本地后，你可以选择将其重命名，移入包含在"$PATH"的目录里并赋予其可执行权限，比如在 Linux 上你可以执行如下命令完成这些操作：

```shell
mv dtm-linux-amd64 /usr/local/bin/dtm
chmod +x dtm
```

接着你可以通过如下命令验证 dtm 的权限以及版本等是否正确：

```shell
$ dtm version
0.9.1
```