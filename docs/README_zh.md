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

**注意：我们会优先更新英文版 [README](../README.md) ，中文版有一定的滞后，强烈建议大家直接阅读英文版**

如果你懒得看下面的一大串文字：一个开源 DevOps 工具链管理工具。

不过我还是建议你看下下面的一大串文字：

假如你现在成立一家公司，或者具体一点，你要组建一个研发团队，在开始写代码前你需要做哪些事情？有些事情是绕不开的，比如：

1. 你需要选择一个地方来存放代码，也许是 GitHub，也许是 GitLab；
2. 你需要一个工具来完成项目管理或者说需求管理、Issue 管理等等工作，也许你会选择 Jira 或者禅道或者 Trello；
3. 你需要选择一种开发语言，选择一个开发框架，比如你决定用 Golang 来开发，假如这是一个 web 项目，你需要考虑 web 框架用什么？“第一行”代码怎么写，也就是第一个脚手架怎么组装；
4. 然后你需要配置一些 ci 自动化，比如 GitHub 上添加 actions 来完成代码的扫描、测试等等；
5. 当然 cd 工具也不能少，不管你选择 Jenkins 还是 ArgoCD；
6. 如果 cd 完成了，接下来可能你马上要开始纠结日志、监控、告警等等方案应该怎么定了
7. 如果想得更多，或许你希望 GitHub 上别人给你提的 issue 能够自动同步到你的 Jira 或者 Trello……
8. ……

也许我上面说到的例子并不完整或者绝对准确，但是有一个结论是我们必须接受的：“在一个软件的开发生命周期中，除了业务代码编码本身，在 DevOps 工具链上我们将花费大量精力去选型、打通、落地、维护……”

所以 DevStream 要解决什么问题呢？我们要做的就是将主流的涵盖 DevOps 全生命周期的开源工具管理起来，包括这些工具的安装部署、最佳实践配置、工具间的打通等等。

## DevStream 目前能干什么？

1. 缺陷、需求管理 - Trello (集成 GitHub)
2. 源码管理 - Golang 脚手架生成
3. CI 流程 - Golang、Python、Nodejs
4. CD/GitOps - ArgoCD / ArgoCD App
5. Monitoring - kube-prometheus
6. ……

## 快速开始

如果你想要快速上手体验，可以跳转到我们的[快速开始](./quickstart_zh.md)文档。

## 你想问 DevStream 的将来？

或许用不了多久，我们就能完整实现 “DevOps toolchain as code”，那时候你的整个 DevOps 工具链都能以 DevStream 作为唯一入口来运维，dtm(DevStream 命令行工具)将成为你的整条 DevOps 工具链的 “single source of truth”。当然那时你需要替换整个 DevOps 工具链中的某一个环节，也会变得很简单。

其实目前我们已经部分实现 “single source of truth”，部署好的工具发生的部分变更已经能够被 dtm 感知到，并且 dtm 会判断这种变更是否合理，是否需要修复，进而采取相应的动作让整个 DevOps 工具链变得更可靠。

## 怎么参与 DevStream 社区？

当然，DevStream 的发展离不开社区用户的支持，DevStream 欢迎所有人参与社区建设，一起完善 dtm 的功能，让 dtm 越来越强大！

不要有任何心理负担，我们非常欢迎大家下载、体验、捉虫、提 Issue、挑刺、bugifx 等等等等。

## 交流、支持

如果你发现了 bug 或者有任何好的意见建议，我们希望你直接在 GitHub 上给我们提 issue。

当然 <a href="https://join.slack.com/t/devstream-io/shared_invite/zt-16tb0iwzr-krcFGYRN7~Vv1suGZjdv4w" target="_blank">Slack</a> 也是有的。

另外我们也有微信群，可以直接扫下方群二维码加入用户群。

![](images/wechat-group-qr-code.jpg)
