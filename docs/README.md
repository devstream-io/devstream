# Documentation

Our docs website is: https://www.devstream.io/docs/index. Read more there.

---

To report a bug or create a new doc, please go to this repo: https://github.com/devstream-io/website.

If you want to know exactly how, read on:

## Creating a New Documentation (English)

First of all, thanks for your effort in improving our documentation!

To create a new document, please do the followings:

### 1 Fork the Website Repo

Fork this repo: https://github.com/devstream-io/website.

### 2 Creating a New Markdown File

In the forked repo, under `docs` folder, create a new Markdown file.

_Note 1: We use number-indexed filenames so that it's easier to understand the order without having to look inside each doc._

 For example:

```shell
.
├── 01-index.md
├── 02-quickstart.md
├── 03-config.md
├── 04-architecture.md
├── 05-core-concepts.md
├── 06-output.md
├── 07-commands-explained.md
├── 08-plugins
├── 09-development-workflow.md
├── 10-project-layout.md
├── 11-creating-a-plugin.md
├── 12-best-practices
└── 13-autocomplete.md
```

If you want to insert a new doc at a certain position, you need to rename those files after it.

_Note 2: in each file, at the very beginning, we have a section like the following:_

```yaml
---
sidebar_position: 1
---
```

This number corresponds to the number index of the file.

Then you can simply write the documentation in Markdown format. 

If you need to insert a picture, the picture should be put under `static/img/docs/`. Refer to existing documents about how to use links to an image file under the static folder.

### 3 i18n

Currently, our doc supports English and Simplified Chinese. English doc is our top priority. For Simplified Chinese, we can put the English doc there without translation.

For example, if you want to create a Simplified Chinese version of the doc `docs/01-index.md`, you can copy it to the folder for Simplified Chinese:

```shell
cp docs/01-index.md i18n/zh/docusaurus-plugin-content-docs/current/
```

Then open the file `i18n/zh/docusaurus-plugin-content-docs/current`, and update the content to Simplified Chinese.

### 4 Create a Pull Request to `devstream-io/website`

---

## 创建一个中文文档

我们的文档以英文为最高优先级，但如果您只想写一份中文文档（并且也不想把它翻译成英文），您可以这样做：

### 1 Fork the Website Repo

Fork我们的文档所在的repo: https://github.com/devstream-io/website.

### 2 创建一个新的Markdown文件

在fork的repo里，在`i18n/zh/docusaurus-plugin-content-docs/current`目录下，
创建一个新的Markdown文件。

_注1: 文件名以数字编号打头，以便我们不用打开文档站也能一眼看出某篇文档的顺序和位置。_

_注2: 在每个文档里最开头部分我们有一段类似于下面的指令：_

```yaml
---
sidebar_position: 1
---
```

这个`sidebar_position`的数字要跟文件名开头的数字一致。

如果你想把你的文件插入在某个特定的位置，你需要让你新创建的文件名和里面的sidebar_position保持一致，并且修改后续的文件。

如果您需要在文档中插入图片，图片请放在`static/img/docs/`目录下。可以参考已有的文档来学习如何引用这个目录下的图片。

### 3 英文文档需要翻译吗？

不需要，但需要放一个占位符。

比如，如果您先创建了`i18n/zh/docusaurus-plugin-content-docs/current/01-index.md`文档，那么请把它复制到英文文档所在的目录下：

```shell
cp i18n/zh/docusaurus-plugin-content-docs/current/ docs/01-index.md
```

然后编辑`docs/01-index.md`文件，删掉英文内容，写一个"Work in Progress"即可，这样DevStream PMC成员就知道这个文档是需要被翻译的了。

### 4 提交PR

您的成果需要以PR的形式提交至devstream-io/website repo.
