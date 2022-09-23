# 安装

## 1 用 curl 安装

进入你的工作目录，运行：

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/install/download.sh)"
```

这个命令会下载 `dtm` 二进制文件到你的工作目录中，并赋予二进制文件执行权限。

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
