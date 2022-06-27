# 构建

```shell
cd path/to/devstream
make clean
make build -j8 # 多线程构建
```

上述命令将构建所有内容: `dtm` 本身及所有插件。

项目也支持以下构建模式：

- 只构建 `dtm`：`make build-core`。
- 构建指定的插件：`make build-plugin.PLUGIN_NAME`。例如： `make build-plugin.argocd`。
- 构建所有插件：`make build-plugins -j8` (多线程构建)。

使用 `make help` 获取更多信息。
