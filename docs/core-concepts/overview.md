# Overview

## 1 Architecture

![](../images/architecture-overview.png)

The diagram above shows how data flows through DevStream modules.

In essence:

- The core of DevStream (plugin engine) acts like a state machine, which calculates actions based on the configuration and state.
- Then DevStream core calls "plugins" to accomplish the CRUD actions of specific DevOps tools or integrations.

---

## 2 Plugin

Plugin is a critical concept of DevStream.

As shown above, DevStream uses a core-plugin architecture where the core acts mainly as a state machine and engine. The core/engine in turn drives the plugins, which are responsible for creating/reading/updating/deleting/integrating DevOps tools.

Plugins are automatically downloaded and managed by dtm according to the config file.

Developers and contributors can write their own plugins to extend the capabilities of DevStream. See [creating a plugin](../development/dev/creating-a-plugin.md) for more detail.

---

## 3 State

The _State_ records the current status of your DevOps toolchain and platform defined and created by DevStream.

The state contains the configuration of all the pieces and their corresponding status so that the DevStream core can rely on it to calculate the required actions to reach the state defined in the config.

---

## 4 Config

DevStream defines desired status of your DevOps platform in config files.

The config can be a single YAML file, as well as a bunch of YAML files under the same directory. The config contains the following sections:

- `config`: basic configuration of DevStream, at the moment mainly state-related settings. Read more [here](./state.md).
- `vars`: variable definitions. Key/value pairs, which can be referred to in the tools/apps/pipelineTemplates sections.
- `tools`: a list of DevStream _Tools_, each containing its name, instanceID (unique identifier), and options. Read more [here](./tools.md).
- `apps`: a list of _Apps_, another DevStream concept, each corresponding to a microservice. Read more [here](./apps.md).
- `pipelineTemplates`: a list of templates which can be referred to by DevStream _Apps_. Read more [here](./apps.md).

---

## 5 Workflow

![config state resource-status workflow](../images/config_state_resource.png)
