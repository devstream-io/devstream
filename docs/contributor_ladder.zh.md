# 贡献者阶梯指南
嗨，你好啊！听说你想要了解 DevStream 项目的贡献者阶梯(contributor ladder) 吗？这太酷了！让我们接着往下看吧！

此贡献者阶梯指南涵盖了 DevStream 项目中不同的贡献者角色及其所对应的责任和特权等。

社区成员一般会从"阶梯"的第一级开始，然后在参与项目中逐步升级。

DevStream 项目的成员会尽可能的帮助所有的贡献者在贡献者阶梯中逐步攀登。

下面是分为三个维度来描述贡献者角色：
- "责任"：贡献者在其级别应该做到的事情。
- "要求"：贡献者在其级别需要满足的资格。
- "特权"：贡献者在其级别被赋予的额外权限。

## 贡献者（Contributor）
_定义：直接为项目做出有价值的贡献。其表现形式可以不只是代码。非代码的贡献同样重要且被我们所认可！在 Contributor 级别的贡献者可能是一个新的 Contributor，或者只是偶尔提交一些贡献。_

- 责任：
  - 遵循 [CNCF CoC](https://github.com/cncf/foundation/blob/main/code-of-conduct.md)
  - 遵循 [项目贡献指南](https://github.com/devstream-io/devstream/blob/main/docs/contributing_guide.zh.md)
- 要求（以下的一项或者多项 PS：多多益善 💖）：
  - 提交 issue 或者偶尔解决一些 issue。
  - 偶尔提交一些 PR。
  - 提交文档类的贡献。
  - 参加 DevStream 社区例会或其他会议。
  - 尽可能解答其他社区成员提出的问题。
  - 提交一些 issue 或 PR 的反馈。
  - 在 releases 时做一些测试工作，如果发现问题可以提交一些补丁类代码，或者帮助 review。
  - 组织或协助组织社区活动。
  - 对 DevStream 项目做一些宣传。
  - 帮助维护项目的基础设施等。
- 特权：
  - 受邀参加 Contributor 活动。
  - 能够进一步成为组织成员（Organization Member）。

## 组织成员（Organization Member）
_定义：组织成员（Organization Member）首先也是贡献者（Contributor），不同的是组织成员（Organization Member）经常性地参与 DevStream 项目。
组织成员（Organization Member）在项目资源库和选举中都拥有特权，他们对项目抱有极大的兴趣。_

组织成员（Organization Member）除了要遵循贡献者（Contributor）的责任，还应遵循以下责任：

- 责任：
  - 定期的做出项目贡献，例如：[Github Insights](https://github.com/devstream-io/devstream/pulse) 所示，每年至少有 50 个 Github 的贡献。
- 要求：
  - 要对项目有提交成功的贡献，至少包括以下一项：
    - 5 个合并的 PR。
    - Review 5 个 PR。
    - 解决并关闭 5 个 issue。
    - 或与以上等效的贡献组合。
  - 参与贡献的时间不低于 1 个月。
  - 至少熟知 1 个项目区域。
  - 得到两名组织成员（Organization Member）推荐，并且这些推荐人不能来自于同一个雇主。
- 特权：
  - 被分配一些 issue 和 review 。
  - 执行 CI/CD 命令。
  - 投票选举权。
  - 加入 @devstream 团队。
  - 推荐其他贡献者成为组织成员（Organization Member）。

贡献者（Contributor）成为组织成员（Organization Member）的流程如下：

1. 贡献者由推荐官在合适的存储库提交 issue 来提名。
2. 来自不同出处的第二位推荐官同意该 issue 的提名。
3. 提名将会在社区会议上公开 review。所有成员对该提名进行投票。需要多数票通过才能批准该提名。
4. 在 issue 中公开宣布结果，且公开宣布在 Slack 的 [devstream-contributors 频道](https://cloud-native.slack.com/messages/devstream-contributors)中。

## 评审者（Reviewer）
_定义：评审者（Reviewer）在项目中有贡献和审查的记录，负责审查特定的代码、文档、测试或其他项目区域。评审者（Reviewer）们对该区域的代码变动的审查集体负责，判断是否可以合并该修改。他们的审查记录将会被作为贡献记录到项目中。 _

评审者（Reviewer）负责"特定区域"。可以是特定的代码目录、驱动程序、文档章节、测试、事件或其他规模小于主库或子项目的组件。多数情况下，可能是 Git 仓库中的一个或一组目录。注:以下的"特定领域"指代的是其负责的特定项目区域。

此外，评审者除了拥有组织成员（Organization Member）的所有权利以外，还有其要肩负的责任。

- 责任：
  - 遵循审查指南。
  - 审查他们负责的特定领域的大部分的 Pull Request。
  - 每年至少审查 20 个 PR。
  - 帮助其他贡献者（Contributor）成为评审者（Reviewer）。

- 要求：
  - 担任贡献者（Contributor）至少满一个月。
  - 是一名组织成员（Organization Member）。
  - 审查或帮助审查过不低于 5 个 Pull Request。
  - 分析并解决了特定领域的测试失败问题。
  - 对某些领域有深入的理解。
  - 承诺负责特定领域。
  - 支持新的和临时的贡献者（Contributor），并帮助贡献者（Contributor）提交对社区有用的 PR 。

- 额外的特权：
  - 拥有 Github 或 CI/CD 的权限，可以批准特定领域的 PR 。
  - 推荐或批准其他的贡献者成为评审者（Reviewer）。
  - 将会作为批准者（Approver），被列在特定领域对应文件夹下的 OWNERS 文件中。

成为评审者的流程：

1. 在对应仓库打开一个 PR ，将提名者的 Github 用户名添加到一个或多个目录的 OWNERS 文件中。
2. 至少要有两名主仓库或该特定领域的批准者（Approver）批准，才可批准此 PR 。

## 维护者（Maintainer）
_定义：维护者（Maintainer）是非常成熟的贡献者，负责整个项目。因此，他们有权批准针项目任何类型的 PR 。同时在有关项目战略和决策上有主导权。_

维护者（Maintainer）满足评审者（Reviewer）所有的责任和要求，并额外满足：

- 责任：每年至少审查 40 个 PR ，尤其是那些跨模块的复杂的 PR 。
  - 指导新加入的评审者（Reviewer）。
  - 撰写重构的 PR 。
  - 参与 CNCF 维护者活动。
  - 确定项目的战略和政策。
  - 参与并主导社区例会。

- 要求：
  - 担任评审者（Reviewer）至少满三个月。
  - 能够对项目的多个领域有深入了解。
  - 能够独立于雇主、朋友或团队，做出利于 DevStream 项目发展的判断。
  - 能够保证每月在 DevStream 项目工作的时间不低于 40 小时。

- 额外的特权：
  - 批准项目中任何的 PR 。
  - 作为维护者（Maintainer）公开代表 DevStream 项目。
  - 代表 DevStream 项目与 CNCF 交流沟通。
  - 拥有维护者决策会议的投票权。

成为维护者（Maintainer）的流程：

1. 任何维护者（Maintainer）都可以提名评审者（Reviewer）成为新的维护者（Maintainer），具体的提名方式是：在 DevStream 仓库的根目录打开一个 PR ，将提名者添加为 OWNERS 文件中的审批者。
2. 提名者需要在上述打开的 PR 中留下评论，确保其承诺遵循维护者（Maintainer）的所有要求。
3. 大多数维护者（Maintainer）持同意观点，并批准该 PR。

## 不活跃措施
任何级别的贡献者们都需要保持活跃，对项目保持兴趣。不活跃有弊于项目，长期不活跃代表贡献者的流失、可能会丧失对项目的信任，导致项目延期。

- 衡量是否为活跃的贡献者，其标准是：
  - 超过两个月没有任何贡献。
  - 超过一个月和社区无任何沟通交流。

- 长期不活跃会导致：
  - 非自愿的撤职或降职。
  - 被要求转为"荣誉退休（Emeritus）"状态。

## 淘汰机制
当贡献者的责任和能力没有满足时，会导致非自愿的移除/降级，可能是由于长时间不活跃或不满足所要求的能力又或者违反了社区的准则。对社区来说，执行淘汰机制很重要，因为它保护了社区及其可交付成果，同时也为新的贡献者提供了加入的机会。

非自愿的移除或降级需要多数维护者（Maintainer）的投票来处理。

## 降级/退出机制
如果贡献者（Contributor）的承诺并没有兑现的话，可以考虑依据贡献阶梯降低，或退出 DevStream 项目。当然我们并不希望有此类事情发生。

请与维护者（Maintainer）们联系，商讨你的降级/退出事宜。

## 联系反馈
如有疑问，请咨询 @PMC主席 @IronCore864.

