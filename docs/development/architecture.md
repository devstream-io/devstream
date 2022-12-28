# Architecture

This document summarizes the main components of DevStream and how data flows between these components.

## 0 Data Flow

The following diagram shows an approximation of how DevStream executes a user command:

![DevStream Architecture Diagram](../../images/architecture-overview.png)

There are three major parts:

- CLI: handles user input, commands and parameters.
- Plugin engine: achieves the core functionalities by calling other modules (config manager, plugin manager, state manager, and backend manager.)
- plugins: implements the actual CRUD interfaces for a certain DevOps tool, or integrates different tools together. Each plugin corresponds to a certain DevOps tool or an automated integration of tools.

## 1 CLI (The `devstream` Package)

Note: for simplicity, the CLI is named `dtm`(DevOps Tool Manager) instead of the full name DevStream.

Every time a user runs the `dtm` program, the execution transfers immediately into one of the "command" implementations in the [`devstream`](https://github.com/devstream-io/devstream/tree/main/cmd/devstream) package, in which folder all commands' definitions reside.

Then, each command calls the plugin engine package under [`internal/pkg`](https://github.com/devstream-io/devstream/tree/main/internal/pkg/pluginengine).

The `pluginengine` calls the [config manager package](https://github.com/devstream-io/devstream/tree/main/internal/pkg/configmanager) first to read the local YAML config file into a struct.

Then it calls the [`pluginmanager` package](https://github.com/devstream-io/devstream/tree/main/internal/pkg/pluginmanager) to download the required plugins.

After that, the plugin engine calls the [state manager](https://github.com/devstream-io/devstream/tree/main/internal/pkg/statemanager) to calculate "changes" between the congfig, the state, and the actual DevOps tool's status. At last, the plugin engine executes actions according to the changes, and updates the state. During the execution, the plugin engine loads each plugin (`*.so` file) and calls the predefined interface according to each change.

## 2 Plugin Engine

The `pluginengine` has various responsibilities:

- make sure the required plugins (according to the config file) are present 
- generate changes according to the config, the state and tools' actual status
- execute the changes by loading each plugin and calling the desired action

It achieves the goal by calling the following modules:

### 2.1 Config Manager

Model types in package [`configmanager`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configmanager/configmanager.go#L23) represent the top-level configuration structure.

### 2.2 Plugin Manager

The [`pluginmanager`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginmanager/manager.go) package is in charge of downloading necessary plugins according to the configuration.

If a plugin with the desired version already exists locally, it will not download it again.

### 2.3 State Manager

The [`statemanager`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/manager.go) package manages the "state", i.e., what has been done successfully and what not.

The `statemanager` stores the state in a [`backend`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/backend/backend.go).

### 2.4 Backend Manager

The [`backend`](https://github.com/devstream-io/devstream/tree/main/internal/pkg/backend) package is the backend manager, which manages the actual state. Currently, local, remote (AWS S3 compatible), and k8s(ConfigMap) state are supported.

## 3 Plugin

A _plugin_ implements the aforementioned, predefined interfaces.

It executes operations like `Create`, `Read`, `Update`, and `Delete`.

To develop a new plugin, see [creating a plugin](dev/creating-a-plugin.md).
