# github-actions 插件

_该插件用于在项目中创建 Github Action Workflows。_

## 用例

下面的内容是一个示例配置文件用于创建 Github Workflows：

``` yaml
--8<-- "github-actions.yaml"
```

**注意:**

- `scm` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。
- `pipeline` 选项项用于控制 `CI` 流程中的各个阶段，具体配置可查看文档[pipline配置项](pipeline.zh.md)。