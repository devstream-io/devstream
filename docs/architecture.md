# DevStream Architecture Summary

This document summarizes the main components of DevStream and how data and requests flow between these components.

## DevStream Request Flow

The following diagram shows an approximation of how DevStream executes a user command:

![DevStream Architecture Diagram](./images/architecture-overview.png)

## CLI (The `devstream` Package)

For simplicity, the CLI is named `dtm` instead of the full name DevStream.

Every time a user runs the `dtm` program, the execution transfers immediately into one of the "command" implementations in the [`devstream`](https://github.com/merico-dev/stream/tree/main/cmd/devstream) package, in which folder all commands' definitions reside. Then, each command calls the corresponding package under [`internal/pkg`](https://github.com/merico-dev/stream/tree/main/internal/pkg).

The flow illustrated above applies to the main DevStream commands like `dtm install`, `dtm uninstall` (_TODO_), and `dtm reinstall` (_TODO_). For these commands, the role of the command is to parse all the config, load each plugin, and call the [predefined interface](https://github.com/merico-dev/stream/blob/main/internal/pkg/plugin/plugin.go#L12).

## Configuration Loader

Model types in package [`config`](https://github.com/merico-dev/stream/blob/main/internal/pkg/config/config.go) represent the top-level configuration structure.

## State Manager

_TODO_

## Plugin Engine

The plugin engine has various responsibilities:

- making sure the required plugins (according to the config) are present (_TODO_)
- generate a plan according to the state (_TODO_)
- execute the plan by loading each plugin and running the desired operation

_Note: there will be a rework of the engine to generate and execute a plan according to the state manager, which hasn't been implemented yet._

## Plugins

A _plugin_ implements the aforementioned predefined interfaces.

It executes operations like install, reinstall (_TODO_), uninstall (_TODO_), and stores state (_TODO_).
