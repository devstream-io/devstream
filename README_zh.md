<div align="center">
<br/>

<img src="docs/images/icon-color.svg" width="120">

# DevStream

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](https://makeapullrequest.com)
![Test](https://github.com/devstream-io/devstream/actions/workflows/main.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/devstream-io/devstream)](https://goreportcard.com/report/github.com/devstream-io/devstream)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/6202/badge)](https://bestpractices.coreinfrastructure.org/projects/6202)
[![Downloads](https://img.shields.io/github/downloads/devstream-io/devstream/total.svg)](https://github.com/devstream-io/devstream/releases)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://cloud-native.slack.com/archives/C03LA2B8K0A)

| [English](README.md) | 中文 |
| --- | --- |

</div>

> We are gonna revive this project recently.

## DevStream 是什么？
TL;DR: DevStream（CLI工具名为`dtm`）是一个开源的DevOps工具链管理器。

[v0.6.0 Demo](https://www.bilibili.com/video/BV1W3411P7oW/)

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

- 我们有很多选择。哪个是最好的？没有"放之四海而皆准"的答案，因为这完全取决于你的需求和喜好。
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

A：受 [`git`](https://github.com/git/git#readme) 的启发，这个名字可以是（取决于你的心情）：

- "**d**evs**t**rea**m**": 一个对称缩写。
- "**D**evops **T**oolchain **M**anager"：当它对你有用的时候。
- "**d**ead **t**o **m**e"：当它崩溃的时候。

## 为什么使用DevStream？

不再需要手动的 `curl/wget` 下载、`apt` 安装、`helm` 安装；不再需要预先的本地试验以保证组件能正确安装。

在一个人类可读的 `YAML` 配置文件中定义你所需要的DevOps工具，只需按一个按钮（或一个命令），你就能建立起整个DevOps工具链和SDLC工作流。

五分钟，一个命令。

想安装另一个不同的工具来试一试？没问题。

想删除或重新安装工作流中的某个特定部分？DevStream已经帮你解决了!

## 快速入门

现在就跟随我们的[快速入门](https://docs.devstream.io/en/latest/quickstart.zh/)文档开始使用 DevStream

## 最佳实践

DevStream支持许多工具的管理。你可以灵活地结合一些工具来满足你所需要的DevOps工具链。

是的，如果你问我是否有可以开箱即用的推荐实践。

我很高兴地告诉你，我们有，而且我们正在不断增加更多可能的组合。

我们非常欢迎你告诉我们你期望的组合。

- [GitOps工具链](https://docs.devstream.io/en/latest/best-practices/gitops.zh/)
- [用 DevStream 搭建 GitLab + Jenkins + Harbor 工具链，管理 Java Spring Boot 项目开发生命周期全流程](https://docs.devstream.io/en/latest/best-practices/gitlab-jenkins-harbor-java-springboot.zh/)

## 支持的DevOps工具

DevStream已经支持许多工具，而且还在不断增加。关于支持的工具的完整列表，请查看我们的 [插件列表](https://docs.devstream.io/en/latest/plugins/plugins-list) 文档。

或者，运行 `dtm list plugins`，它将显示所有可用的插件。

## 开发指南

### 前提条件

- Git
- Go (1.18版本以上)

### 开发指南

- [开发环境搭建](https://docs.devstream.io/en/latest/development/dev/dev-env-setup.zh)
- [代码格式检查](https://docs.devstream.io/en/latest/development/dev/lint.zh)
- [源码构建](https://docs.devstream.io/en/latest/development/dev/build.zh)
- [代码测试：单元测试(unit test)、端到端测试(e2e test)](https://docs.devstream.io/en/latest/development/dev/test.zh)
- [开发新插件](https://docs.devstream.io/en/latest/development/dev/creating-a-plugin.zh)

## 贡献

首先，感谢你愿意为DevStream做贡献！

关于如何贡献、贡献者成长计划、风格指南等更多细节，请查看我们的 [CONTRIBUTING](CONTRIBUTING.md) 文档。

## 社区

我们将定期组织 "DevStream Community Meeting"，请访问 [WIKI](https://github.com/devstream-io/devstream/wiki) 页面了解详情。

另外，请加入我们的Slack channel，请参见如下步骤：

1. [给自己发送进入CNCF的Slack的邀请](https://slack.cncf.io)。
    - 输入邮箱地址，点击"get my invite"。
    - 打开收件箱，找到邀请邮件，点击邮件中的"join now"。
    - 你可以用邮箱地址或者Google账号等方式加入，跟随屏幕提示即可。
2. 对于会说英文的用户：加入DevStream的英语频道。以下两种方式均可：
    - [点击这里](https://cloud-native.slack.com/messages/devstream)。
    - 在Slack app里，在左侧的导航栏，找到“频道”区域，将鼠标移动到上面，这时“频道”的右侧会出现一个加号。点击加号，然后选择“浏览频道”。输入"devstream"，找到频道然后加入。
3. 对于只说中文的用户，可以加入[devstream-mandarin](https://cloud-native.slack.com/messages/devstream-mandarin)频道，这里的讨论都是用中文进行的。

对于使用微信的用户，请扫描二维码进群：

![](docs/images/wechat-group-qr-code.png)

## 行为守则

[DevStream行为守则](/CODE_OF_CONDUCT.md)

DevStream于2022年6月加入CNCF沙盒。我们也需要遵循[云原生计算基金会（CNCF）社区行为准则](https://github.com/cncf/foundation/blob/main/code-of-conduct-languages/zh.md).
