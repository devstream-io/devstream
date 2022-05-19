# Tool, State, and Resource

The architecture documentation explains how in general DevStream works. If you haven't read it yet, make sure you do that before continuing with this document.

## 1 Tool

- One of the major part of the _Config_ is a list of tools, defined in [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/config.go#L23).
- Each _Tool_ has its Name, Plugin, and Options, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/config.go#L40).
- Each _Tool_ can have its dependencies, which are specified by the `dependsOn` keyword.

The dependency `dependsOn` is an array of strings, with each element being a dependency. Each dependency is named in the format of "NAME.PLUGIN". See [here](https://github.com/devstream-io/devstream/blob/main/examples/tools-quickstart.yaml#L12) for example.

## 2 State

- The _State_ is actually a map of states, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L24).
- Each state in the map is a struct containing Name, Plugin, Options, and Resource, as defined [here](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L16).

## 3 Resource

- We call what the plugin created a _Resource_, and the `Read()` interface of that plugin returns a description of that resource, which is in turn stored as part of the state.

Config-State-Resource workflow:

![config state resource workflow](../images/config_state_resource.png)
