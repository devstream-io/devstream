# 为DevStream创建中文文档

## 背景

1. 我们使用[readthedocs](https://readthedocs.org/)托管文档，不过您并不需要对readthedocs了解很多，即可为DevStream的文档做出贡献。
1. Readthedocs支持[`sphinx`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-sphinx.html)和[`mkdocs`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-mkdocs.html)来构建文档站点；我们选择了`mkdocs`。如果您遇到任何问题，请参阅[`mkdocs`的官方文档](https://www.mkdocs.org/)。
1. 我们使用了`mkdocs`的`material`主题。有关`material`主题的更多信息请参见：
    - [网站](https://squidfunk.github.io/mkdocs-material/)
    - [文档](https://squidfunk.github.io/mkdocs-material/getting-started/)
    - [GitHub 仓库](https://github.com/squidfunk/mkdocs-material)
1. 我们还使用了`mkdocs`的两个插件：
    - [search](https://squidfunk.github.io/mkdocs-material/setup/setting-up-site-search/)
    - [mkdocs-i18n](https://pypi.org/project/mkdocs-i18n/)

## 前置条件

- Python3（`mkdocs`是基于Python的）
- pip3（安装依赖项）

建议macOS用户进行如下操作：

- [使用Brew安装Python](https://docs.brew.sh/Homebrew-and-Python)
- `pip3 install -r docs/requirements.txt`

## 文档根目录

`mkdocs`的根文件夹位于`/`。

主配置是根目录`/`下的`mkdocs.yml`，docs文件夹是`/docs`。

## i18n（国际化）

目前，我们支持两种语言：
- 英文
- 简体中文

值得注意的是，搜索功能不适用于简体中文（`mkdocs`搜索功能的限制）。

对于每个英文文档，在单独的文件中都有中文翻译。 如果英文文档的文件名是`doc_name.md`，那么应该还有一个名为`doc_name.zh.md`的文件。要创建中文文档，请将中文翻译内容放入`doc_name.zh.md`。该文件是`doc_name.md`（英文）的翻译。

## 创建一个新文档

要创建新文档，请执行以下操作：

- 在`/docs`文件夹中创建`doc_name.md`和`doc_name.zh.md`。如有必要，您可以将它们放在恰当的子文件夹下。参考当前目录结构来确定适合该文档的最佳路径。
- 编写文档的内容。您可以选择只写英文文档或中文文档；您不必用两种语言编写文档；当然如果您希望展示您中英双语的实力的话，我们建议您两种语言的文档同时编写，一并提交。
- 更新`/mkdocs.yml`文件，更新`nav:`部分。这是整个文档网站的目录。

## 在本地查看您的更改

运行：

```sh
# 在repo的根目录下
pip3 install -r docs/requirements.txt
mkdocs serve
```

请在提交PR之前在本地确认您的更改。

## 推荐的工具和阅读材料

- [Markdown指南](https://www.markdownguide.org/)
- [Grammarly](https://app.grammarly.com/)（用于编写英文文档）
- [标题大小写转换](https://www.titlecase.com/)
