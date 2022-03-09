# Contributor Ladder Growth Programs

**Note:** This document is wip and welcomes everyone to help improve.

## DevStream Community Membership

| Role        | Responsibilities                        | Requirements                                                 | Defined by                                                   |
| ----------- | --------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Contributor | N/A                                     | At least 1 contribution to the project                       | N/A                                                          |
| Member      | Active contributor in the community     | Sponsored by 2 reviewers and multiple contributions to the project | DevStream GitHub Org member & Team member                    |
| Reviewer    | Review contributions from other members | Highly experienced active member who is knowledgeable about the codebase; Member for at least 1 month | DevStream GitHub Org member & Team member  & OWNERS file reviewer entry |
| Approver    | Contributions acceptance approval       | Highly experienced active reviewer and contributor; Reviewer for at least 3 months | DevStream GitHub Org member & Team member & OWNERS file approver entry |

New contributors should be welcomed to the DevStream community by existing members, helped with PR workflow, and directed to relevant documentation and communication channels.

## OWNERS file

`OWNERS` files are used to designate responsibility for different parts of the DevStream codebase. We use them to assign the `Reviewer` and `Approver` roles.

We will gradually define `OWNERS` for each DevStream plugin and each module of DevStream core.

A typical OWNERS file looks like below:

```yaml
approvers:
- daniel # GitHub username
reviewers:
- daniel
- danny
```
