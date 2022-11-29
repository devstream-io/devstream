# Overview

## 1 Architecture

![](../images/architecture-overview.png)

The diagram above shows how data flows through DevStream modules.

In essence:

- The core of DevStream (plugin engine) acts like a state machine, which calculates actions based on the configuration and state.
- Then DevStream core calls "plugins" to accomplish the CRUD actions of specific DevOps tools or integrations.

---

## 2 State

The _State_ records the current status of your DevOps toolchain and platform defined and created by DevStream.

The state contains the configuration of all the pieces and their corresponding status so that the DevStream core can rely on it to calculate the required actions to reach the state defined in the config.

---

## 3 Config

DevStream defines desired status of your DevOps platform in config files.

The main config file, which defaults to `config.yaml` in the working directory, defines where to store the DevStream state, where to load DevStream plugins and the location of other config files.

There are a few different configs, but please don't be overwhelmed because some are not mandatory, and [you can define all things within a single file](https://stackoverflow.com/questions/50788277/why-3-dashes-hyphen-in-yaml-file).

Configurations in the main config contains multiple sections:

- `config`: basic configuration of DevStream, at the moment mainly state-related settings. Read more [here](./state.md).
- `vars`: variable definitions. Key/value pairs, which can be referred to in the tools/apps/pipelineTemplates sections.
- `tools`: a list of DevStream _Tools_, each containing its name, instanceID (unique identifier), and options. Read more [here](./tools-apps.md).
- `apps`: a list of _Apps_, another DevStream concept, each corresponding to a microservice. Read more [here](./tools-apps.md).
- `pipelineTemplates`: a list of templates which can be referred to by DevStream _Apps_. Read more [here](./tools-apps.md).

---

## 4 Workflow

![config state resource-status workflow](../images/config_state_resource.png)
