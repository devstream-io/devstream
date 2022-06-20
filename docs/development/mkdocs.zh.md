# 为DevStream创建中文文档

## 背景

1. 我们使用[readthedocs](https://readthedocs.org/)托管文档，不过您并不需要对readthedocs了解很多即可为DevStream创建文档。
2. Readthedocs支持[`sphinx`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-sphinx.html)和[`mkdocs`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-mkdocs.html)来构建文档站点；我们现在使用`mkdocs`。如果您遇到任何问题，`mkdocs`的官方文档[这里](https://www.mkdocs.org/)是一个很好的起点。
3. 在我们的mkdocs设置中，我们使用`material`主题。有关”material”的更多信息：
    - [网站](https://squidfunk.github.io/mkdocs-material/)
    - [文档](https://squidfunk.github.io/mkdocs-material/getting-started/)
    - [GitHub 仓库](https://github.com/squidfunk/mkdocs-material)
4. 我们还为mkdocs使用了两个插件：
    - [search](https://squidfunk.github.io/mkdocs-material/setup/setting-up-site-search/)
    - [mkdocs-i18n](https://pypi.org/project/mkdocs-i18n/)

1,860 / 5,000
翻译结果
## 前提条件

- Python3（`mkdocs`是Python实现的）
- pip3（安装依赖项）

macOS用户的建议设置：

- [使用brew设置Python](https://docs.brew.sh/Homebrew-and-Python)
-`pip3 install -r docs/requirements.txt`

## 文档根目录

`mkdocs`的根文件夹位于`/`。

主配置是根目录`/`下的`mkdocs.yml`，docs文件夹是`/docs`。

## i18n（国际化）

目前，我们支持两种语言：
- en
- zh

值得注意的是“search”功能不适用于“zh”（这是因为“mkdocs”搜索功能的限制）

对于每个英文文档，在单独的文件中都有中文翻译。

如果英文文档的文件名是`doc_name.md`，那么应该还有一个名为`doc_name.zh.md`的文件。要创建中文文档，请将中文翻译内容放入`doc_name.zh.md`。该文件是`doc_name.md`（英文）的翻译。

## 创建一个新文档

要创建新文档，请执行以下操作：

- 在`/docs`文件夹中创建`doc_name.md`和`doc_name.zh.md`。如有必要，您可以将它们放在子文件夹下。参考当前目录结构并使用您的最佳判断来确定该新文档的最佳位置。
- 编写文档的内容。您可以选择只写英文文档或中文文档；您不必（但如果可以的话，强烈推荐）用两种语言编写文档。
- 更新`/mkdocs.yml`文件，更新`nav:`部分。这是整个文档网站的目录。

## 在本地查看您的更改

运行以下命令：

```sh
# 在repo的根目录
pip3 install -r docs/requirements.txt
mkdocs serve
```

在提交PR之前检查您的更改。

## 推荐参考资源

- [Markdown指南](https://www.markdownguide.org/)
- [Grammarly](https://app.grammarly.com/)（用于编写英文文档）
- [标题大小写转换](https://www.titlecase.com/)
