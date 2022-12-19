<div align="center">
<br/>

<img src="./docs/images/icon-color.svg" width="120">

# DevStream

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](https://makeapullrequest.com)
![Test](https://github.com/devstream-io/devstream/actions/workflows/main.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/devstream-io/devstream)](https://goreportcard.com/report/github.com/devstream-io/devstream)
[![Downloads](https://img.shields.io/github/downloads/devstream-io/devstream/total.svg)](https://github.com/devstream-io/devstream/releases)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/6202/badge)](https://bestpractices.coreinfrastructure.org/projects/6202)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://cloud-native.slack.com/archives/C03LA2B8K0A)

| English | [中文](README_zh.md) |
| --- | --- |

</div>

## DevStream, What Is It Anyway?

TL;DR: DevStream (CLI tool named `dtm`) is an open-source DevOps toolchain manager.

[v0.6.0 Demo](https://www.youtube.com/watch?v=q7TK3vFr1kg)

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

## Quick Start

If you want to get a quick start, follow our [quick start](https://docs.devstream.io/en/latest/quickstart/) doc now.

## Best Practices Toolchain Integration

DevStream supports the management of many tools. You can flexibly combine some tools to meet the DevOps toolchain your need.

And yes, if you ask me if any recommended practices that can be used out of the box,

I am happy to tell you that we have, and we are constantly adding more possible combinations,

so you are more than welcome to tell us what combinations you expect.

- [GitOps Toolchain](https://docs.devstream.io/en/latest/best-practices/gitops/)
- [GitLab, Jenkins and Harbor On Premise Toolchain (Chinese only for now)](https://docs.devstream.io/en/latest/best-practices/gitlab-jenkins-harbor-java-springboot.zh/)

## Supported DevOps Tools

DevStream already supports many tools and it's still growing. For a complete list of supported tools, check out our [list of plugins](https://docs.devstream.io/en/latest/plugins/plugins-list/) document.

Alternatively, run `dtm list plugins` and it will show you all the available plugins.

## Dev Info

### Pre-requisites

- Git
- Go (1.18+)

### Development Guide

- [Development Environment Setup](https://docs.devstream.io/en/latest/development/dev/dev-env-setup)
- [Code linter](https://docs.devstream.io/en/latest/development/dev/lint)
- [Build the source code](https://docs.devstream.io/en/latest/development/dev/build)
- [Test the source code: unit test, e2e test](https://docs.devstream.io/en/latest/development/dev/test)
- [Create a plugin](https://docs.devstream.io/en/latest/development/dev/creating-a-plugin)

## Contribute

First of all, thanks for wanting to contribute to DevStream! For more details on how to contribute, contributor growth program, style guide and more, please check out our [CONTRIBUTING](./CONTRIBUTING.md) document.

## Community

We will regularly organize `DevStream Community Meeting`, please visit the [wiki](https://github.com/devstream-io/devstream/wiki) page for details.

Please join our Slack channel. Here's how:

1. [Invite yourself to CNCF's Slack if you haven't done so](https://slack.cncf.io). 
    - Input your email address, and click "get my invite."
    - Open your inbox, find the invitation email, and click "join now."
    - You can join by Email or with your Google account; follow the instructions.
2. Join DevStream channel, there are two ways to do so:
    - Use [this link](https://cloud-native.slack.com/messages/devstream) to join the channel.
    - In your Slack app, on the left side navigation bar, move your mouse to the "Channels" section, and there should emerge a "plus" sign on the right. Click, and select "browse channels." Input "devstream", and join.
3. For Mandarin-speaking users and contributors, you are also encouraged to join the [devstream-mandarin](https://cloud-native.slack.com/messages/devstream-mandarin) channel, where all discussions will be in Mandarin.

For WeChat users, you can also join our WeChat group:

![](docs/images/wechat-group-qr-code.png)

## Code of Conduct

[DevStream code of conduct](./CODE_OF_CONDUCT.md)

As of Jun 2022, we joined CNCF sandbox. We also need to follow the [CNCF Community Code of Conduct](https://github.com/cncf/foundation/blob/main/code-of-conduct.md).
