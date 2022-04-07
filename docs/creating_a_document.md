# Creating a Document

## Dependencies

```shell
brew install python
pip3 install sphinx sphinx-rtd-theme
cd docs
pip3 install -r requirements.txt
```

## Locally Build Docs and View

```shell
cd docs
rm -rf build
make html
open build/html/index.html
```

## Where to Put a Document

### Top Level Document

Put it under docs/*.md, if it's a top-level document.

## Non Top Level Document

Put it under:
- docs/best_practices/*.md, if it's a "best practice" type of documentation
- docs/commands_in_detail/*.md, if it's related to `dtm` commands
- docs/plugins/*.md, if it's a doc of a plugin

## Table of Content

### Top Level Document

In `index.md`, add the name of the file (without .md) into the "toctree" part (order matters).

### Non Top Level Document

There should be a top level document for non-top level documents, and of course that top level document should be in the table of content in `index.md`.

In this top level document, there should be a sub table of content.

See `plugins.md` and `plugins/*.md` for an example.

## Document Titles

Always use H1 (#) as the title of the document. There should only be one H1 (#) per doc.

Use H2 (##), H3 (###), etc., accordingly.

## Add the Doc to the Table of Content

Always put the following at the end of your doc:

````yaml
```{toctree}
---
maxdepth: 1
---
```
````