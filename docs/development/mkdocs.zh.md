# 为DevStream创建中文文档

## 背景

1. 我们使用[readthedocs](https://readthedocs.org/)托管文档，不过你并不需要对readthedocs了解很多，即可为DevStream的文档做出贡献。
1. Readthedocs支持[`sphinx`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-sphinx.html)和[`mkdocs`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-mkdocs.html)来构建文档站点；我们选择了`mkdocs`。如果你遇到任何问题，请参阅[`mkdocs`的官方文档](https://www.mkdocs.org/)。
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

- 在`/docs`文件夹中创建`doc_name.md`和`doc_name.zh.md`。如有必要，你可以将它们放在恰当的子文件夹下。参考当前目录结构来确定适合该文档的最佳路径。
- 编写文档的内容。你可以选择只写英文文档或中文文档；你不必用两种语言编写文档；当然如果你希望展示你中英双语的实力的话，我们建议你两种语言的文档同时编写，一并提交。
- 大多数情况下，你不需要考虑导航菜单，它是整个文档网站的目录。但是如果需要自定义导航菜单，可以参考[设置导航](https://github.com/devstream-io/devstream/blob/main/docs/development/mkdocs.md#setting-up-navigation)。

## 设置导航

如果要自定义导航菜单，可以更新`mkdocs.yaml`中的`nav:`部分，支持通配符和子目录链接。例如：

```
nav:
  - DTM Commands Explained in Depth:
    - commands/autocomplete*.md
    - commands/*.md
  - Plugins:
    - plugins/plugins-list*.md
    - plugins/*.md
  - Best Practices: best-practices/
  - 'contributing_guide*.md'
  - 'contributor_ladder*.md'
```

- 通常，`contributing_guide*.md`指代了`contributing_guide.md`和`contributing_guide.zh.md`两个文件，分别是英文和中文文档；
- 如果在`commands/`, `plugins/`, `best-practices/`目录中创建文档，不需要更新`nav`。

如果想了解更多关于导航的配置，请参阅[配置页面与导航](https://www.mkdocs.org/user-guide/writing-your-docs/#configure-pages-and-navigation)和[导航语法](https://oprypin.github.io/mkdocs-literate-nav/)。

## 在本地查看你的更改

运行：

```sh
# 在repo的根目录下
pip3 install -r docs/requirements.txt
mkdocs serve
```

请在提交PR之前在本地确认你的更改。

## 推荐的工具和阅读材料

- [Markdown指南](https://www.markdownguide.org/)
- [Grammarly](https://app.grammarly.com/)（用于编写英文文档）
- [标题大小写转换](https://www.titlecase.com/)
