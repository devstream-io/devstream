# 自动补全

## Fig 自动补全

为了给终端用户提供更好的体验，我们支持[Fig 自动补全](https://github.com/withfig/autocomplete/blob/master/src/dtm.ts)。

与其他自动补全工具不同，[Fig](https://fig.io) 更直观。它为终端用户带来了 IDE 风格的体验。详细介绍见[官方网站](https://fig.io/)

![](fig/fig-intro.gif)

**注意：Fig 目前仅支持 MacOS！**

### 安装 Fig

详见 [https://fig.io](https://fig.io)

![](fig/fig-terminal.png)

安装完成后，你需要将 Fig 集成到你正在使用的终端。

### 补全示例

#### 查看插件模板配置时，补全插件名

![](fig/cmd-show-plugins.gif)

#### 提示子命令

![](fig/cmd-help.gif)

#### 构建时，补全插件名

![](fig/cmd-make.gif)

## Shell 自动补全

### Bash 自动补全

#### Linux

**注：主要参考了 [Linux 上的 bash 自动补全](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#enable-shell-autocompletion)**

自动补全脚本依赖于 `bash-completion`, 需要提前安装：

```bash
apt-get install bash-completion # Ubuntu

yum install bash-completion # CentOS 和 RedHat
```

上面的命令会创建 `/usr/share/bash-completion/bash_completion`，这是 `bash-completion` 的主要脚本。取决于不同的包管理器，你需要手动在 `~/.bashrc` 中引入(`source`) 这个文件。

为了测试是否配置成功，请重新加载 shell ，并运行 `type _init_completion` 来确认。若命令执行成功，则配置正确，否则请在 `~/.bashrc` 中添加如下内容：

```bash
source /usr/share/bash-completion/bash_completion
```

重新加载 shell ，并运行 `type _init_completion` 来确认是否安装、配置成功。

然后，你可以通过 `dtm completion bash` 命令生成 Bash 自动补全脚本，执行下述命令，以将相应内容添加到 `~/.bashrc` 中：

```bash
echo 'source <(dtm completion bash)' >>~/.bashrc
source ~/.bashrc
```

在重新加载 shell 后，dtm 自动补全应该就可以正常使用了！

#### MacOS

!!! note

    主要参考了 [MacOS 上的 bash 自动补全](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/#enable-shell-autocompletion)


自动补全脚本依赖于 `bash-completion`, 需要提前安装：

```bash
brew install bash-completion@2
```

正如以上命令的输出所述，请在 `~/.bash_profile` 文件中添加以下内容：

```bash
brew_etc="$(brew --prefix)/etc"
echo "[[ -r \"${brew_etc}/profile.d/bash_completion.sh\" ]] && . \"${brew_etc}/profile.d/bash_completion.sh\"" >>~/.bash_profile
source ~/.bash_profile
```

重新加载 shell，用 `type _init_completion` 验证 bash-completion v2 是否正确安装。

接着，用 `dtm completion bash` 命令为 Bash 生成自动补全脚本。 执行以下命令，以将相应内容添加到 `~/.bashrc` 中：

```bash
echo 'source <(dtm completion bash)' >>~/.bash_profile
source ~/.bash_profile
```

在重新加载 shell 后，dtm 自动补全应该就可以正常使用了！

### Zsh 自动补全

你可以用 `dtm completion zsh` 命令为 Zsh 生成自动补全脚本。 执行以下命令，以将相应内容添加到 `~/.zshrc` 中：

```zsh
echo 'source <(dtm completion zsh)' >>~/.zshrc
source ~/.zshrc
```

在重新加载 shell 后，dtm 自动补全应该就可以正常使用了！

### Fish 自动补全

你可以用 `dtm completion fish` 命令为 Fish 生成自动补全脚本。然后将以下行添加到 `~/.config/fish/config.fish` 文件中：

```fish
dtm completion fish | source
```

在重新加载 shell 后，dtm 自动补全应该就可以正常使用了！

### PowerShell 自动补全

你可以用 `dtm completion powershell` 命令为 PowerShell 生成自动补全脚本。然后将以下行添加到 `$PROFILE` 文件中：

```powershell
dtm completion powershell | Out-String | Invoke-Expression
```

这个命令会在每次 PowerShell 启动时重新生成自动补全脚本。你也可以直接将生成的脚本添加到 `$PROFILE` 文件中。

为了将生成的脚本添加到 `$PROFILE` 文件中，请在 PowerShell 提示符中运行以下行：

```powershell
dtm completion powershell >> $PROFILE
```

在重新加载 shell 后，dtm 自动补全应该就可以正常使用了！
