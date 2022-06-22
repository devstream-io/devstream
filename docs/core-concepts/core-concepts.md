# Core Concepts Overview

Let's assume you're familiar with core Git, Docker, Kubernetes, Continuous Integratoin, Continuous Delivery, and GitOps concepts. 
These are the core concepts of DevStream.

## The Architecture of DevStream

The architecture image below shows how DevStream works in general:
![](./images/architecture-overview.png)

## Workflow

![config state resource workflow](../images/config_state_resource.png)

## Config, Tool, State and Resource

The architecture documentation explains how in general DevStream works. If you haven't read it yet, make sure you do that before continuing with this document.

### 1. Config

DevStream defines your DevOps toolchain in config files.

The major part of the config file is definition of your DevOps _Tools_.

Each _Tool_ has its name, instanceID (unique identifier), and options for that tool.

Each _Tool_ can also have its dependencies, and can refer to other tool's output as values of its own options.

### 2. Tool

//TODO(daniel-hutao): @IronCore Please help to update the descriptions below.

- One of the major part of the _Config_ is a list of tools, defined in [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/config.go#L23).
- Each _Tool_ has its Name, InstanceID, and Options, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/config.go#L40).
- Each _Tool_ can have its dependencies, which are specified by the `dependsOn` keyword.

The dependency `dependsOn` is an array of strings, with each element being a dependency. Each dependency is named in the format of "TOOL_NAME.INSTANCE_ID". 
See [here](https://github.com/devstream-io/devstream/blob/main/examples/quickstart.yaml#L22) for example.

### 3. State

The _State_ records the current status of your DevOps toolchain. It contains the configuration of each tool, and the current status.

- The _State_ is actually a map of states, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L24).
- Each state in the map is a struct containing Name, Plugin, Options, and Resource, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L16).

### 4. Resource

- We call what the plugin created a _Resource_, and the `Read()` interface of that plugin returns a description of that resource, which is in turn stored as part of the state.

Config-State-Resource workflow:

![config state resource workflow](../images/config_state_resource.png)
