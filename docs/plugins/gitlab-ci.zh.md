# gitlab-ci 插件

该插件用于在使用 gitlab 的项目仓库中创建 Gitlab CI 并运行对应项目的 gitlab runner。

### 配置项
下面的内容是一个示例配置文件用于创建 Gitlab CI：

``` yaml
--8<-- "gitlab-ci.yaml"
```

**注意:**

- `scm` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。
- `pipeline` 选项项用于控制 `CI` 流程中的各个阶段，具体配置可查看文档[pipline配置项](pipeline.zh.md)。
