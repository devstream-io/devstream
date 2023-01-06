# 概览

## 这篇文章是为"我"准备的吗？

- 你是否正准备创建一个 web 应用/后端 API 微服务？
- 你是否使用 GitHub 存储你的源代码？
- 你是否渴望快速而自动地在你的开发/测试环境中部署你的服务，一触即发？

如果你回答"是"，那么"GitHub Actions + Argo CD + Python + Flask" 工具集可解君愁。

请听我细细道来。

---

## 这篇文章真的是为"我"准备的吗？

这个 "GitHub + GitHub Actions + Argo CD + Python + Flask" 的组合（也可以是 "GitLab + GitLab CI"），如果安装/管理/集成得当，可能会为你构建最佳的 DevOps 平台，如若：

- 你选择 GitHub 作为你的源代码管理（SCM）系统。它可以是：
    - GitHub cloud，或
    - GitHub Enterprise（本地）
- 由于你使用 GitHub 作为 SCM，而 CI 与代码仓库交互频繁，因此使用 GitHub Actions 作为你的 CI 选择是方便的。
    - 当然，你可以选择其他 CI，如 Travis CI、CircleCI 等，它们需要与 GitHub 进行一些额外集成操作；而使用 GitHub Actions，实际上几乎不需要集成。
    - 此外，因为用的人多了，你遇到的问题大概率别人已经遇到过，所以你可以轻松地通过简单的 Google 搜索找到解决方案。
- 你希望在开发环境中快速且自动地部署你的服务。
    - GitOps 可能是你的最佳选择，因为它快速、自动化。
    - 你的环境不多，这可以避免一些棘手的问题，比如如果你使用 GitOps 部署到多个环境，你可能会遇到 Version propagation(多环境下的版本发布，例如将某个通过测试环境的版本发布到生产环境）的问题。
- 你想构建一些后端 API 类型的应用。
    - 根据一份流行度排名，Python 现在是世界上最流行的编程语言。当然，编程语言的流行度会随时间而起伏。但根据 [TIOBE](https://www.tiobe.com/tiobe-index/)，一家荷兰的软件质量保证公司，它一直在跟踪编程语言的流行度，指出“在过去 20 年中，我们第一次看到了一种新的领导者：Python 编程语言。Java 和 C 的长期霸权已经结束。”虽然 Python 可能不是每种情况下的最佳选择，但当你无法掌握全部的信息，且必须做出一个猜测时，选择最流行的语言也不是最糟糕的选择。
    - Flask 是一个小型且轻量级的 Python Web 框架，它提供了很多有用的工具，包括快速调试器、内置的开发服务器以及更多让 Python 创建 Web 应用程序变得简单的特性。

---

## 能做什么？怎么做？

我们想要构建一个基于这些工具的 DevOps 平台，而这个平台可以做到非常复杂的事情：

一个 "宏图" 胜过千言万语：

![](../../images/gitops-workflow.png)

如果你还是更喜欢文字而不是图片（我知道有些人是这样的），我便来讲讲这中间发生的故事（DevStream 为实现这一流程而帮你做的事情）：

- 在你的 GitHub 账户中创建一个仓库，并生成一些代码（一个 Python Flask 应用程序，包含 Dockerfile、helm chart 等）。
- 创建持续集成（CI）流水线，于是：
    - 当 pull request 被创建时，会自动触发一个任务，以运行单元测试。
    - 当 pull request 被合并到主分支时，另一个任务将被触发，来构建 Docker 镜像，并推送到 Dockerhub，再触发部署流程。
- 安装 Argo CD 以进行持续部署（CD）。
- 以 GitOps 的方式触发 CD：当 pull request 被合并到主分支时，它将应用程序部署到 Kubernetes 集群。

让我们"大展宏图"吧！

---

## 建业之基（先决条件）

1. 一个可用的 Kubernetes 集群，用于部署 GitOps 所需的工具，以及作为开发环境。
2. 一个 GitHub 账户和一个用于 API 访问的访问令牌，以及一个 Dockerhub 账户和一个用于 API 访问的 access token。
