# Creating a Documentation for DevStream

## Background

1. We use [readthedocs](https://readthedocs.org/) to host our documentation. You do not need to know too much more about it in order to create a doc for DevStream, though.
1. Readthedocs supports both [`sphinx`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-sphinx.html) and [`mkdocs`](https://docs.readthedocs.io/en/stable/intro/getting-started-with-mkdocs.html) to build the doc site; we use `mkdocs` at the moment. If you meet any issues, the official doc of `mkdocs` [here](https://www.mkdocs.org/) is a good place to start.
1. In our mkdocs setup, we use the `material` theme. More info on "material":
    - [website](https://squidfunk.github.io/mkdocs-material/)
    - [docs](https://squidfunk.github.io/mkdocs-material/getting-started/)
    - [GitHub repo](https://github.com/squidfunk/mkdocs-material)
1. We also use two plugins for mkdocs:
    - [search](https://squidfunk.github.io/mkdocs-material/setup/setting-up-site-search/)
    - [mkdocs-i18n](https://pypi.org/project/mkdocs-i18n/)

## Prerequisites

- Python3 (`mkdocs` is a Python thing)
- pip3 (to install dependencies)

Suggested setup for macOS users:

- [use brew to setup Python](https://docs.brew.sh/Homebrew-and-Python)
- `pip3 install -r docs/requirements.txt`

## Docs Root

The root folder of `mkdocs` is at `/`.

The main config is `/mkdocs.yml`, and the docs folder is `/docs`.

## i18n (Internationalization)

Currently, we support two languages:
- en
- zh

It's worth noting that the "search" function doesn't work for "zh" (a limitation of `mkdocs`' search function.)

For each English document, there is a Chinese translation in a separate file.

If the English document's filename is `doc_name.md`, there should also be a file named `doc_name.zh.md`. To create a Chinese translation, put the translation into `doc_name.zh.md`. This file is the translation of `doc_name.md` (English).

## Create a New Documentation

To create new documentation, do the following:

- Create `doc_name.md` and `doc_name.zh.md` in the `/docs` folder. You can put them under a subfolder if necessary. Refer to the current directory structure and use your best judgment to decide the best place for that new doc.
- Write the content of the doc. You can choose to write only the English doc or the Chinese doc; you don't have to (but it's highly recommended if you can) write documentation in both languages.
- Update the `/mkdocs.yml` file, update the `nav:` section. It's the table of content of the whole doc website.

## Review Your Change Locally

Run:

```sh
# at the root of the repo
pip3 install -r docs/requirements.txt
mkdocs serve
```

Review your changes before submitting a PR.

## Recommended Resources

- [Markdown Guide](https://www.markdownguide.org/)
- [Grammarly](https://app.grammarly.com/) (for writting English documents)
- [Title Case Converter](https://www.titlecase.com/)
