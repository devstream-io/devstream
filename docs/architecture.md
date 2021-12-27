# DevStream Architecture Overview

This document summarizes the main components of DevStream and how data flows between these components.

## 0. Data Flow

The following diagram shows an approximation of how DevStream executes a user command:

![DevStream Architecture Diagram](./images/architecture-overview.png)

There are three major parts:

- CLI: handles user input, commands and parameters
- plugin engine: achieves the core functionalities by calling other modules (config loader, plugins manager, etc.)
- plugins: implements the actual install/uninstall for a certain DevOps tool

## 1. CLI (The `devstream` Package)

Note: for simplicity, the CLI is named `dtm`(DevOps Tool Manager) instead of the full name DevStream.

Every time a user runs the `dtm` program, the execution transfers immediately into one of the "command" implementations in the [`devstream`](https://github.com/merico-dev/stream/tree/main/cmd/devstream) package, in which folder all commands' definitions reside.

Then, each command calls the plugin engine package under [`internal/pkg`](https://github.com/merico-dev/stream/tree/main/internal/pkg/pluginengine).

The plugin engine calls the config loader first to read the local YAML config file into a struct.

Then it calls the plugin manager to download the required plugins.

After that, the plugin engine calls the state manager and plan manager to create a plan. At last, the plugin engine asks the plan manager to execute the plan and the state manager to update the state. During the execution of the plan, the plugin engine loads each plugin (*.so file) and calls the predefined interface.

## 2. Plugin Engine

The plugin engine has various responsibilities:

- making sure the required plugins (according to the config) are present 
- generate a plan according to the state
- execute the plan by loading each plugin and running the desired operation

It achieves the goal by calling the following modules:

### 2.1 Configuration Loader

Model types in package [`config`](https://github.com/merico-dev/stream/blob/main/internal/pkg/configloader/config.go#L12) represent the top-level configuration structure.

### 2.2 Plugin Manager

The [plugin manager](https://github.com/merico-dev/stream/blob/main/internal/pkg/pluginmanager/manager.go) is in charge of downloading necessary plugins according to the configuration.

If a plugin with the desired version already exists locally, it will not download it again.

### 2.3 State Manager

The [state manager](https://github.com/merico-dev/stream/blob/main/internal/pkg/statemanager/manager.go) manages the "state", i.e., what has been done successfully and what not.

### 2.4 Plan Manager

The [plan manager](https://github.com/merico-dev/stream/blob/main/internal/pkg/planmanager/plan.go) creates a plan according to the state and the config and executes the plan.

## 3. Plugins

A _plugin_ implements the aforementioned predefined interfaces.

It executes operations like install, uninstall, and stores state.

To develop a new plugin, see [creating_a_plugin.md](https://github.com/merico-dev/stream/blob/main/docs/creating_a_plugin.md).
