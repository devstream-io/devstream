# 贡献指引

## 致新贡献者
欢迎你的到来，非常感谢你愿意一起建设 DevStream 💖。

初次参与，你遇到任何问题都可以直接反馈与联系社区，包括但不限于：

- 开发环境搭建时遇到任何问题。
- 快速开始(quick start) 文档遇到任何问题。
- 自动化脚本的任何问题。

如果你在运行项目的时候发现任何不符合预期或不合理的地方，请直接提交 Bug 报告！

## 如何贡献
我们欢迎各种贡献，包括但不限于：

- 新功能（feature）
- 代码构建、CI/CD
- Bug 修复
- 文档
- Issue 分类、发起、回复、管理、维护等（Issue Triage）
- 在 微信群、Sack、邮件列表解答问题
- 网站页面设计
- 在各种媒体、博客文章宣传 DevStream
- 发版管理

不是只有提交 Pull Request 才能与 DevStream 擦出火花，你还可以参与我们的 [社区会议](https://github.com/devstream-io/devstream/wiki)，或者直接联系我们([slack](https://cloud-native.slack.com/archives/C03LA2B8K0A)、[微信群](https://raw.githubusercontent.com/devstream-io/devstream/main/docs/images/wechat-group-qr-code.png)）一起聊聊怎么共建社区。

## 参与社区会议
DevStream 真诚地欢迎每一个人参与我们的会议，不需要被邀请，直接来就行！哪怕你暂时没有贡献的想法，只是来听听，我们也非常热烈欢迎。来听，就够了。

- 你可以在我们的 [wiki 页面](https://github.com/devstream-io/devstream/wiki)) 找到关于社区会议的详细信息
- 也请加入我们的 [Slack 频道](https://cloud-native.slack.com/archives/C03LA2B8K0A) ，因为会议日程和重要的事情会在 Slack 公布。
- 对于中国用户，也可以加 [微信群](https://raw.githubusercontent.com/devstream-io/devstream/main/docs/images/wechat-group-qr-code.png)。
- 第一次参会你可以不打开摄像头，简单地介绍自己就足够了。当然，如果你想多聊几句，无论是 Q&A 还是聊天，我们都非常欢迎。
- 随着时间的推移与你对 DevStream 沟通的加深，我们希望你能自如地发表意见、对他人的想法提出反馈，甚至分享你的想法与经验。

## 寻找 Issue
你可以在 Issue 列表找到这样两个标签： `good first issue` 表示仅向新贡献者开放，`help wanted` 适合于所有贡献者。

- [good first issue](https://github.com/devstream-io/devstream/labels/good%20first%20issue) 含有更多的描述信息与指导，往往只涉及一小部分，完成它不需要你熟悉整个项目。如果你是第一次参与 DevStream（甚至是第一次参与开源），非常推荐你从 good first issue 开始共建之旅。更多信息，请参阅 [good first issue 文档](./development/good-first-issues.zh.md)。
- [help wanted](https://github.com/devstream-io/devstream/labels/help%20wanted) 指引了你在完成了 good first issue 后，适合继续贡献的地方。带有这个标签的 issue，除了核心贡献者外的任何人都可参与。
- 你不需要拘泥于带有这两个标签的 issue，其他任何你感兴趣的 issue，都可以直接参与，无论是方案修改建议、讨论与回复、还是代码。
- 有时 good first issue 或 help wanted 因为社区过于热情导致暂时没有空余，只要你愿意，你仍可以参与贡献！直接在 [Slack](https://join.slack.com/t/devstream-io/shared_invite/zt-16tb0iwzr-krcFGYRN7~Vv1suGZjdv4) 或 [微信群](https://raw.githubusercontent.com/devstream-io/devstream/main/docs/images/wechat-group-qr-code.png) 联系我们，告诉我们你愿意参与！

若你想负责某个 issue，请在 issue 内留下评论，例如："I want to work on this"。

## 寻求帮助
在贡献时，向我们提问的最好的方式如下：

- 直接在 GitHub issue 下评论
- [Slack 频道](https://cloud-native.slack.com/archives/C03LA2B8K0A)
- [微信群](https://raw.githubusercontent.com/devstream-io/devstream/main/docs/images/wechat-group-qr-code.png) 也可以

我们更推荐你在 GitHub 或 Slack 提问，这有助于技术讨论被沉淀下来，以帮助后来的贡献者。

## Pull Request 流程与指南
针对 PR、贡献者和审阅者的不成文规则（我们不想有那么多条条框框，但这确实能为你减少不必要的时间浪费）。

- 我们鼓励贡献者提交 PR，即使它处于“work-in-progress”状态。
- 如果你的 PR 代码还没完成编写，请使用 GitHub 的 "[draft PR](https://github.blog/2019-02-14-introducing-draft-pull-requests/)" 功能来告诉社区这个 PR 还在编写中。同时，这也意味着，你不一定要在所有代码完成后再 push，可以先实现大部分内容，提交个 draft PR，这有助于与社区及时交流沟通。
- 当你的 PR 做好了被审阅（review）(处于非 draft PR 状态）的准备后，审阅者应在 24 小时内进行初始审阅。如果是周末或节假日，可能会延迟。
- PR 的作者应当主动请求审阅（request-review），或者在 PR 内评论与 @审阅者，也可直接通过社交媒体进行沟通。如果审阅者未在24小时内回复，请相信我们绝不是有意遗漏你，请再次提醒我们。
- 微小的修改（甚至是拼写错误）也完全可以提 PR。不存在所谓的"小 PR"，不必担心 PR 过小导致代码不被合并。
- 无论是较大的 PR 还是较小的，我们建议你创建一个对应的 issue 以持续跟进与追踪。
- 我们没有特定的特性分支。通常，你会将 PR 并入主分支。如果它针对于特定某个版本的错误修复，请确保你的 PR 请求合并的是该特定分支。
- 如果你想贡献，但不想通过提 PR 来贡献，可以在 Slack 频道、微信群。或直接在 Issue、PR 下告诉我们，以便 maintainer 直接提交 PR。
- 如果你长期未回复 maintainer（通常是一个月），你的 PR 可能会被关闭，敬请谅解。
- 目前，我们不会定期发布，因此无法保证你的 PR 何时包含在下一个版本中。但是，我们正在尽最大努力使发布尽可能频繁。
- 为了鼓励贡献者，你可以在未100%完成原有设计的情况下，完成相应的任务，即关闭 issue 与合并 PR。在这种情况下，你可以再创建个新的 issue，并提交个新的 PR 来完成之前遗漏的工作。我们不希望你在一个 PR 上耗时太久，这会打击你的兴趣。

## 开发环境搭建
- 代码格式检查：[golangci-lint](https://github.com/golangci/golangci-lint)
- 推荐的 IDE：[Visual Studio Code](https://code.visualstudio.com/), [GoLand](https://www.jetbrains.com/go/)
- [文档](https://docs.devstream.io/en/latest/index.zh/)
- [快速开始](https://docs.devstream.io/en/latest/quickstart.zh/)
- [源码](https://github.com/devstream-io/devstream)
- [源码构建](https://docs.devstream.io/en/latest/development/build.zh/)
- [代码测试：单元测试(unit test)、端到端测试(e2e test)](https://docs.devstream.io/en/latest/development/test.zh/)
- TODO: 源码测试, 集成/端到端测试
- TODO: 生成本地文档预览

## 代码提交署名（Sign Off）
授权与认证（licensing) 对于开源项目非常重要，它确保了软件能基于作者提供的条款继续运作。我们需要你在贡献代码的时候署名你的提交，[Developer Certificate of Origin (DCO)](https://developercertificate.org/) 是一种认证你编写了此段代码，并表明你持有这段代码的方式。

你可以通过在提交信息（Git Commit Message）中附加这段信息。注意，你的署名必须与 Git 的用户名和邮箱对应。

    This is my commit message

    Signed-off-by: Your Name <your.name@example.com>

Git 有一个 `-s` 的命令行选项可以帮助你自动署名:

    git commit -s -m 'This is my commit message'

如果你忘记了署名操作，但未将代码提交到远程仓库，可以通过这行命令来补充署名信息：

    git commit --amend -s 

## Pull Request 检查项清单

当你提交 PR 或向其推送新提交时，我们的自动化系统将对你的新代码运行一些检查。我们要求你的 PR 通过这些检查，但在我们接受和合并它之前，我们还有更多的标准。我们建议你在提交代码之前在本地检查以下事项：

- 检查代码格式问题([golangci-lint](https://github.com/golangci/golangci-lint)): 虽然我们的 CI 会在提交代码时自动运行 lint，但在本地先执行 lint 确保无误后再提交能节省你更多的精力与时间。
- 构建/测试 代码，同上。
- 仔细检查你的提交信息（Commit Message）：确认它们符合 [约定式提交](https://www.conventionalcommits.org/zh-hans/v1.0.0/) 规范。同样的，CI 也会执行这一检查，但若你在提交代码前先行检查能加快 PR 合并速度。

## Maintainer 团队

DevStream 由最初由 [思码逸 Merico](https://www.crunchbase.com/organization/merico) 的工程师们发起，并由 [@IronCore864](https://github.com/ironcore864) 带领。

目前，maintainer 们为：

- @IronCore864
- @daniel-hutao

> 注意：虽然 DevStream 由 Merico 发起，以及现在的 maintainer 们都来自 Merico，但我们不希望开源社区中存在 _任何_ 专制的行为。我们将促进长期参与 member 和 reviewer 们加入 maintainer 队伍。我们欢迎来自任何公司与组织的你成为我们社区的一员。

你可以阅读 [贡献者成长阶梯计划](https://docs.devstream.io/en/latest/contributor_ladder.zh/) 以了解详情。
