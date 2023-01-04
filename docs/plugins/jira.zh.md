# jira 插件

该插件用于将 Github 中的 Issue 实时同步到 Jira 中。

## 用例

```yaml
--8<-- "jira.yaml"
```

**注意:**

- Jira 项目的语言必须是英语。
- Jira 的项目必须是已存在。
- `scm` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。
- `jira.token` 需要先在 Jira 中创建，如何创建可以查询文档 [Manage API tokens for your Atlassian account](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/)。
