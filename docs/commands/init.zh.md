# dtm init

`dtm` 将插件发布在 AWS S3 Bucket。

当运行 `dtm init` 时，它会从 AWS S3 下载插件(通过Cloudflare CDN）。

`dtm` 提供了两种下载插件的方式：

1. 根据配置文件：根据配置中的 `tools` 和 `apps` 下载插件
2. 根据命令行参数：根据命令行的参数下载插件

## 根据配置文件下载插件

这种方式下，`dtm init` 将根据配置中定义的 `tools` 和 `apps` 下载所需的插件。

**命令：** `dtm init -f <config file/config dir>`。

你可以把所有的配置放在一个文件中，也可以把配置分散至同一个目录下的多个以 `yaml` 或 `yaml` 为后缀的文件中。

关于配置文件、tools and apps，详见[DevStream 配置](../core-concepts/config.zh.md)。

## 根据命令行参数下载插件

这可以用来预先下载插件，以在**离线环境**使用 `dtm`。

**命令：**

- 下载指定插件，如：`dtm init --download-only --plugins="repo-scaffolding, github-actions" -d=.devstream/plugins`。 
- 下载所有插件，如：`dtm init --download-only --all -d=.devstream/plugins`。


## Init 逻辑

- 根据配置文件和 tool file，或命令行参数，决定下载哪些插件。
- 如果该插件在本地存在，而且版本正确，就什么都不做。
- 如果缺少该插件，则下载该插件。
- 下载后，`dtm` 还会校验该插件的md5值。

## 命令行参数

| 短              | 长            | 默认                      | 描述                                              |
|-----------------|---------------|--------------------------|--------------------------------------------------|
| -f              | --config-file | `"config.yaml"`          | 配置文件路径                                       |
| -d              | --plugin-dir  | `"~/.devstream/plugins"` | 插件存储目录                                       |
| --download-only | `false`       | `""`                     | 只下载插件，一般用于离线使用 `dtm`                    |
| -p              | --plugins     | `""`                     | 要下载的插件，用","隔开，应与 --download-only 一起使用 |
| -a              | --all         | `false`                  | 下载所有插件，应与 --download-only 一起使用          |
|                 | --os          | 本机的操作系统             | 要下载插件的操作系统。                               |
|                 | --arch        | 本机的架构                 | 要下载插件的架构。                                  |



