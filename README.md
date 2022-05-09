<div align="center">
<br/>

![](./docs/images/logo-120px.jpg)

# DevStream

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
![Test](https://github.com/devstream-io/devstream/actions/workflows/main.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/devstream-io/devstream)](https://goreportcard.com/report/github.com/devstream-io/devstream)
[![Downloads](https://img.shields.io/github/downloads/devstream-io/devstream/total.svg)](https://github.com/devstream-io/devstream/releases)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devstream-io/shared_invite/zt-16tb0iwzr-krcFGYRN7~Vv1suGZjdv4w)

| English | [中文](docs/README_zh.md) |
| --- | --- |

</div>

## DevStream, What Is It Anyway?

TL;DR: DevStream (CLI tool named `dtm`) is an open-source DevOps toolchain manager.

Imagine you are starting a new project or ramping up a new team. Before writing the first line of code, you have to figure out the tools to run an effective SDLC process and from development to deployment.

Typically, you'd need the following pieces in place to work effectively:

- Project management software or issue tracking tools (JIRA, etc.)
- Source code management (GitHub, Bitbucket, etc.)
- Continuous integration tools (Jenkins, CircleCI, Travis CI, etc.)
- Continuous delivery/deployment tools (Flux CD/Flux2, Argo CD, etc.)
- A single source of truth for secrets and credentials (secrets manager, e.g., Vault by HashiCorp)
- Some tools for centralized logging and monitoring (for example, ELK, Prometheus/Grafana);

The list could go on for quite a bit, but you get the idea!

There are many challenges in creating an effective and personalized workflow:

- There are too many choices. Which is best? There is no "one-size-fits-all" answer because it totally depends on your needs and preferences.
- Integration between different pieces is challenging, creating silos and fragmentation.
- The software world evolves fast. What's best today might not make sense tomorrow. If you want to switch parts or tools out, it can be challenging and resource intensive to manage.

To be fair, there are a few integrated products out there that may contain everything you might need, but they might not suit your specific requirements perfectly. So, the chances are, you will still want to go out and do your research, find the best pieces, and integrate them yourself. That being said, to choose, launch, connect, and manage all these pieces take a lot of time and energy.

You might be seeing where we are going with this...

We wanted to make it easy to set up these personalized and flexible toolchains, so we built DevStream, an open-source DevOps toolchain manager.

Think of the Linux kernel V.S. different distributions. Different distros offer different packages so that you can always choose the best for your need.

Or, think of `yum`, `apt`, or `apk`. You can easily set it up with your favorite packages for any new environment using these package managers.

DevStream aims to be the package manager for DevOps tools.

To be more ambitious, DevStream wants to be the Linux kernel, around which different distros can be created with various components so that you can always have the best components for each part of your SDLC workflow.

## Why `dtm`?

Q: The CLI tool is named `dtm`, while the tool itself is called DevStream. What the heck?! Where is the consistency?

A: Inspired by [`git`](https://github.com/git/git#readme), the name is (depending on your mood):

- a symmetric, scientific acronym of **d**evs**t**rea**m**.
- "devops toolchain manager": you're in a good mood, and it actually works for you.
- "dead to me": when it breaks.

## Why Use DevStream?

No more manual curl/wget download, apt install, helm install; no more local experiments and playing around just to get a piece of tool installed correctly.

Define your desired DevOps tools in a single human-readable YAML config file, and at the press of a button (one single command), you will have your whole DevOps toolchain and SDLC workflow set up. Five Minutes. One Command.

Want to install another different tool for a try? No problem.

Want to remove or reinstall a specific piece in the workflow? DevStream has got your back!

## Installation

Please visit GitHub [Releases](https://github.com/devstream-io/devstream/releases) page and download the appropriate binary according to your operating system and architecture.

## Quick Start

If you want to get a quick start, follow our [quick start](./docs/quickstart_en.md) doc now.

## Configuration

This is an example of DevStream config: [examples/tools-quickstart.yaml](./examples/tools-quickstart.yaml).

Remember to open this configuration file, modify all FULL_UPPER_CASE_STRINGS (like YOUR_GITHUB_USERNAME, for example) in it to your own.

Pay attention to the meaning of each item to ensure that it is what you want.

For other plugins, checkout the "Plugins" section in our [doc](https://www.devstream.io/docs/index) for detailed usage.

## Usage

To apply the config, run:

```shell
./dtm apply -f YOUR_CONFIG_FILE.yaml
```

If you don't specify the config file with the "-f" parameter, it will try to use the default value which is "config.yaml" from the current directory.

_`dtm` will compare the config, the state, and the resources to decide whether a "create", "update", or "delete" is needed. For more information, read our [Core Concepts documentation here](https://www.devstream.io/docs/core-concepts)._

The command above will ask you for confirmation before actually executing the changes. To apply without confirmation (like `apt-get -y update`), run:

```shell
./dtm -y apply -f YOUR_CONFIG_FILE.yaml
```

To delete everything defined in the config, run:

```shell
./dtm delete -f YOUR_CONFIG_FILE.yaml
```

_Note that this deletes everything defined in the config. If some config is deleted after apply (state has it but config not), `dtm delete` won't delete it. It differs from `dtm destroy`._

Similarly, to delete without confirmation:

```shell
./dtm -y delete -f YOUR_CONFIG_FILE.yaml
```
To delete everything defined in the config, regardless of the state:

```shell
./dtm delete --force -f YOUR_CONFIG_FILE.yaml
```

To verify, run:

```shell
./dtm verify -f YOUR_CONFIG_FILE.yaml
```

To destroy everything, run:

```shell
./dtm destroy
```

_`dtm` will read the state, then determine which tools are installed, and then remove those tools. It's same as `dtm apply -f empty.yaml` (empty.yaml is an empty config file)._

## Best Practices Toolchain Integration

DevStream supports the management of many tools. You can flexibly combine some tools to meet the DevOps toolchain your need.

And yes, if you ask me if any recommended practices that can be used out of the box,

I am happy to tell you that we have, and we are constantly adding more possible combinations,

so you are more than welcome to tell us what combinations you expect.

- [GitOps Toolchain](https://www.devstream.io/docs/best-practices/gitops)

## Supported DevOps Tools

DevStream already supports many tools and it's still growing. For a complete list of supported tools, check out our [list of plugins](https://www.devstream.io/docs/plugins/plugins-list) document.

Alternatively, run `dtm list plugins` and it will show you all the available plugins.

## Dev Info

### Pre-requisites

- Git
- Go (1.17+)

### Build

```shell
cd path/to/devstream
make clean
make build -j8 # multi-threaded build
```

This builds everything: `dtm` and all the plugins.

We also support the following build modes:
- Build `dtm` only: `make build-core`.
- Build a specific plugin: `make build-plugin.PLUGIN_NAME`. Example: `make build-plugin.argocd`.
- Build all plugins: `make build-plugins -j8` (multi-threaded build.)

See `make help` for more information.

### Test

Run all unit tests:

```shell
go test ./...
```

e2e test runs on GitHub actions.

## Contribute

First of all, thanks for wanting to contribute to DevStream! For more details on how to contribute, contributor growth program, style guide and more, please check out our [CONTRIBUTING](./CONTRIBUTING.md) document.

## Community

We will regularly organize `DevStream Community Meeting`, please visit the [wiki](https://github.com/devstream-io/devstream/wiki) page for details.

- Message us on <a href="https://join.slack.com/t/devstream-io/shared_invite/zt-16tb0iwzr-krcFGYRN7~Vv1suGZjdv4w" target="_blank">Slack.</a>
- For Chinese users, the WeChat group QR code is as below:

![](docs/images/wechat-group-qr-code.jpg)
