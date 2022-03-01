# Config, State, Resource Detailed Explanation and How We Manage Changes

The architecture documentation explains how in general DevStream works. If you haven't read it yet, make sure you do that before continuing with this document.

## 1. Config, State, and Resource

Config
- The _Config_ is a list of tools, defined in [here](https://github.com/merico-dev/stream/blob/main/internal/pkg/configloader/config.go#L19).
- Each _Tool_ has its Name, Plugin, and Options, as defined [here](https://github.com/merico-dev/stream/blob/main/internal/pkg/configloader/config.go#L24).

State

- The _State_ is actually a map of states, as defined [here](https://github.com/merico-dev/stream/blob/main/internal/pkg/statemanager/state.go#L21).
- Each state in the map is a struct containing Name, Plugin, Options and Resource, as defined [here](https://github.com/merico-dev/stream/blob/main/internal/pkg/statemanager/state.go#L14).

Resource
- We call what the plugin created a _Resource_, and the `Read()` interface of that plugin returns a description of that resource, which is in turn stored as part of the state.

## 2. Changes for `dtm apply`

When _applying_ a config file using `dtm`, here's what happens:

- Read the _State_
- For each _Tool_ defined in the _Config_, we compare the _Tool_, its _State_, and the _Resoruce_ it has created before (if the state exists). We create some changes based on that.
- For each _State_ that doesn't have a _Tool_ in the _Config_, we generate a "Delete" change to delete the _Resource_. Since there isn't a _Tool_ in the config but there is a _State_, it means maybe the _Resource_ had been created previously then the user removed the _Tool_ from the _Config_, which means the user don't want the _Resource_ any more.

## 3. Changes for `dtm delete`

When _deleting_ using `dtm`, here's what happens:

- Read the _Config_ only
- For each _Tool_ defined in the _Config_, if there is a corresponding _State_, we generate a "Delete" change.

## 4. Changes for `dtm destroy`

TODO.
