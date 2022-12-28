# dtm develop

`dtm develop` is used to develop a new plugin for anything.

## 1 Brief

`dtm develop create-plugin` generates plugin scaffolding based on dtm [pre-built templates](https://github.com/devstream-io/devstream/tree/main/internal/pkg/develop/plugin/template).

`dtm develop validate-plugin` is used for verifying whether all the necessary files required for the plugin exist

## 2 Full Process Guide

A detailed guide to the plugin development process can be found in [creating a plugin](../development/dev/creating-a-plugin.md).

## 3 Flags

`dtm develop create-plugin`:

| Short | Long   | Default | Description                              |
|-------|--------|---------|------------------------------------------|
| -n    | --name | `""`    | specify name of the plugin to be created |

`dtm develop validate-plugin`:

| Short | Long   | Default | Description                                |
|-------|--------|---------|--------------------------------------------|
| -n    | --name | `""`    | specify name of the plugin to be validated |
| -a    | --all  | `false` | validate all plugins                       |
