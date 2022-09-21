# Install

## 1 Install dtm binary with curl

In your working directory, run:

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/install/download.sh)"
```

This will download the `dtm` binary to your working directory, and grant the binary execution permission.

> Optional: you can then move `dtm` to a place which is in your PATH. For example: `mv dtm /usr/local/bin/`.

## 2 Install with [asdf](https://asdf-vm.com/)

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
