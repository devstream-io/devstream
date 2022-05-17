<div align="center">
<br/>

![](./images/logo-120px.jpg)

# DevStream

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
![Test](https://github.com/devstream-io/devstream/actions/workflows/main.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/devstream-io/devstream)](https://goreportcard.com/report/github.com/devstream-io/devstream)
[![Downloads](https://img.shields.io/github/downloads/devstream-io/devstream/total.svg)](https://github.com/devstream-io/devstream/releases)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devstream-io/shared_invite/zt-16tb0iwzr-krcFGYRN7~Vv1suGZjdv4w)

| [English](../README.md) | 中文 |
| --- | --- |

</div>

## DevStream 是什么？
TL;DR: DevStream（CLI工具名为`dtm`）是一个开源的DevOps工具链管理器。

想象你正在开始一个新的项目或组建一个新的团队。在写第一行代码之前，你需要一个能够高效运转SDLC(软件开发生命周期)和承载开发至部署全过程的工具。

通常情况下，你需要以下几个部分来高效地工作。

- 项目管理软件或 `issue` 追溯工具（JIRA等）
- 源代码管理（GitHub、Bitbucket等）
- 持续集成（Jenkins、CircleCI、Travis CI等）
- 持续交付/部署（Flux CD/Flux2、Argo CD等)
- 密钥和证书的单一事实来源(A single source of truth)（密钥管理器，如HashiCorp的Vault）
- 集成化的日志和监控工具（例如，ELK、Prometheus/Grafana）
- ......

具体内容远远不止这些，不过你应该已经明白意思了。

在创建一个高效、定制化的工作流上，当前有许多挑战。

- 我们很多选择。哪个是最好的？没有"放之四海而皆准"的答案，因为这完全取决于你的需求和喜好。
- 不同部分之间的整合是非常具有挑战性的，否则将导致项目孤岛化、碎片化。
- 软件领域演进很快。今天最好的东西可能明天就毫无意义。如果你想换掉一些组件或工具，管理起来会很困难，也很耗费资源。

说实话，有一些产品可能包含你需要的一切，但它们可能并不完全适合你的具体要求。因此，你仍然需要自己去搜寻，找到最好的组件，并自己将它们整合起来。也就是说，选择、启动、连接和管理所有这些组件需要大量的时间和精力。

你可能已经看到了我们想要做的事情......

我们想简化整合组件的过程，所以我们建立了DevStream，一个开源的DevOps工具链管理器。

想一想Linux内核与不同发行版的关系。不同的发行版提供不同的软件包，这样你就可以随时选择你最需要的。

或者，想想`yum`、`apt`或`apk`。你可以使用这些包管理器为任何新环境轻松设置你最喜欢的软件包。

**DevStream的目标是成为DevOps工具的软件包管理器。**

**更具野心的是，DevStream想成为Linux内核，你可以用各种组件创建不同的发行版，为SDLC工作流的每个部分选择最适合的组件。**

## 为什么是 `dtm` ？
Q：CLI被命名为 `dtm`，而工具本身被称为 `DevStream`。这是怎么回事！？一致性在哪里？

A：受 [`git`](https://github.com/git/git#readme)的启发，这个名字可以是（取决于你的心情）：

- "**d**evs**t**rea**m**": 一个对称缩写。
- "**D**evops **T**oolchain **M**anager"：你的心情很好，而且它确实对你有用。
- "**d**ead **t**o **m**e"：当它崩溃的时候。

## 为什么使用DevStream？

不再需要手动的 `curl/wget` 下载、`apt` 安装、`helm` 安装；不再需要预先的本地试验以保证组件正确安装。

在一个人类可读的 `YAML` 配置文件中定义你所需要的DevOps工具，只需按一个按钮（或一个命令），你就能建立起整个DevOps工具链和SDLC工作流。

五分钟，一个命令。

想安装另一个不同的工具来试一试？没问题。

想删除或重新安装工作流中的某个特定部分？DevStream已经帮你解决了!
