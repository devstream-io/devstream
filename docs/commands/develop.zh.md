# dtm develop

`dtm develop` 用于开发一个新的插件，做任何你想做的事情。

## 1 简介

- `dtm develop create-plugin` 根据[dtm 插件模板](https://github.com/devstream-io/devstream/tree/main/internal/pkg/develop/plugin/template)生成插件脚手架。
- `dtm develop validate-plugin` 用于验证插件所需的必要文件是否都存在。

## 2 完整插件开发指南

插件开发的详细指南可以在[创建一个插件](../development/dev/creating-a-plugin.zh.md)找到。

## 3 命令行参数

`dtm develop create-plugin`:

| 短  | 长     | 默认值 | 描述             |
|-----|--------|-------|-----------------|
| -n  | --name | `""`  | 要创建的插件的名称 |

`dtm develop validate-plugin`:

| 短  | 长     | 默认值   | 描述             |
|-----|--------|---------|-----------------|
| -n  | --name | `""`    | 要验证的插件的名称 |
| -a  | --all  | `false` | 验证所有插件      |
