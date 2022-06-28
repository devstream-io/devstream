# Contributing Guide

## New Contributor Guide

Welcome! We are glad that you want to contribute to DevStream! 💖

As you get started, you are in the best position to give us feedback on areas of our project that we need help with include:

- problems found during setting up a new developer environment;
- gaps in our quick start guide or documentation;
- bugs in our automation scripts.

If anything doesn't work when you run it or doesn't make sense, please open a bug report and let us know!

## Ways to Contribute

We welcome many different types of contributions including:

- New features
- Builds, CI/CD
- Bug fixes
- Documentation
- Issue Triage
- Answering questions on Slack/Mailing List
- Web design
- Communications / Social Media / Blog Posts
- Release management

Not everything happens through a GitHub pull request. Please come to our [meetings](https://github.com/devstream-io/devstream/wiki) or [contact us](https://cloud-native.slack.com/archives/C03LA2B8K0A) and let's discuss how we can work together. 

## Come to Meetings

Absolutely everyone is welcome to come to any of our meetings. You never need an invitation to join us. We want you to join us, even if you don’t have anything you feel like you want to contribute. Just being there is enough!

- You can find out more about our meetings on [our wiki page here](https://github.com/devstream-io/devstream/wiki).
- Please also consider joining our [Slack channel](https://cloud-native.slack.com/archives/C03LA2B8K0A) because meeting schedules and announcements will also be made available there.
- You don’t have to turn on your video. The first time you come, introducing yourself is more than enough.
- Over time, we hope that you feel comfortable voicing your opinions, giving feedback on others’ ideas, and even share your ideas and experiences.

## Find an Issue

We have good first issues for new contributors and help wanted issues suitable for any contributor.

- [good first issue](https://github.com/devstream-io/devstream/labels/good%20first%20issue) has extra information to help you make your first contribution. If you are new to DevStream (or even new to open-source,) this is a good place to get yourself started. For more information, see the [good first issues doc](./development/good-first-issues.md).
- [help wanted](https://github.com/devstream-io/devstream/labels/help%20wanted) are issues suitable for someone who isn't a core maintainer and is good to move onto after your "good first issue."
- Sometimes there won’t be any issues with these labels. That’s ok! There is likely still something for you to work on. If you want to contribute but you don’t know where to start or can't find a suitable issue, you can reach out to us on our [Slack channel](https://cloud-native.slack.com/archives/C03LA2B8K0A) and ask for an issue to work on.

Once you see an issue that you'd like to work on, please post a comment saying that you want to work on it. Something like "I want to work on this" is fine.

## Ask for Help

The best way to reach us with a question when contributing is to ask:

- The original GitHub issue
- Our [Slack channel](https://cloud-native.slack.com/archives/C03LA2B8K0A).

## Pull Request Lifecycle

The unwritten rules (heck, we are writing them anyways) for PRs, contributors, and reviewers:

- Contributors are encouraged to submit a PR even when it's in a work-in-progress state.
- If a PR is still a work-in-progress, use the "[draft PR](https://github.blog/2019-02-14-introducing-draft-pull-requests/)" function of GitHub to indicate you are already working on it, and show the community members your progress.
- When a PR is ready for review (not in the draft PR status,) reviewers are expected to do an initial review within 24 hours. If it's a weekend or holiday, it might get delayed.
- The author should ping/bump when the pull request is ready for further review. If the PR appears to be stalled for over 24 hours, you can also ping/bump.
- If your PR is stuck (seems can't get reviewed,) besides ping/bump, you can also reach out to the [Slack channel](https://cloud-native.slack.com/archives/C03LA2B8K0A) and notify the reviewers there.
- Small scope (incremental value, even one spelling correction) PR is completely fine. No PR is too small. Don't worry that because you only changed a word and that PR won't get merged.
- Complete feature PRs are also welcome. In this case, however, try not to make the PR so huge that it becomes extremely hard for reviewers to review.
- Small or big PR, it's nice to create an issue for each PR to keep things tracked.
- We don't have a specific feature branch. Normally you would do a PR against the main branch. If it's a bugfix specific to a version, make sure a PR to that release branch is also created.
- If you as a contributor do not want to follow through with a PR, please try to let us know via our Slack channel or directly in the PR/issue, so that the maintainers can continue working on it by committing to the PR directly to help get it merged.
- Maintainers might close a PR if the contributor hasn’t responded in a month.
- Currently, we don't release regularly, so there isn't any guarantee about when your PR will be included in the next release. However, we are trying our best to make releases as frequent as possible.
- As an encouragement to contributors, it's fine to close an issue or merge a PR even if the originally designed features aren't 100% implemented. In cases like this, please encourage contributors to create a follow-up issue and PR to further implement that. Do not make the PR process last too long because it might discourage contributors.

## Development Environment Setup

- linter: [golangci-lint](https://github.com/golangci/golangci-lint)
- recommended IDE: [Visual Studio Code](https://code.visualstudio.com/), [GoLand](https://www.jetbrains.com/go/).
- [docs](https://docs.devstream.io/en/latest/)
- [quick start](./quickstart.md)
- Get the source code: https://github.com/devstream-io/devstream
- [Build the source code](./development/build.md)
- [Test the source code, unit](./development/test.md)
- TODO: Test the source code, integration/end-to-end
- TODO: Generate and preview the documentation locally

## Sign Your Commits

Licensing is important to open source projects. It provides some assurances that the software will continue to be available based on the terms that the author(s) desired. We require that contributors sign off on commits submitted to our project's repositories. The [Developer Certificate of Origin (DCO)](https://developercertificate.org/) is a way to certify that you wrote and have the right to contribute the code you are submitting to the project.

You sign off by adding the following to your commit messages. Your sign-off must match the git user and email associated with the commit.

    This is my commit message

    Signed-off-by: Your Name <your.name@example.com>

Git has a `-s` command-line option to do this automatically:

    git commit -s -m 'This is my commit message'

If you forgot to do this and have not yet pushed your changes to the remote repository, you can amend your commit with the sign-off by running 

    git commit --amend -s 

## Pull Request Checklist

When you submit your pull request, or you push new commits to it, our automated systems will run some checks on your new code. We require that your pull request passes these checks, but we also have more criteria than just that before we can accept and merge it. We recommend that you check the following things locally before you submit your code:

- Lint your code. Although this will be checked by our CI, linting it yourself locally first before creating a PR will save you some time and effort.
- Build/test your code. Same as above.
- Double-check your commit messages. See if they meet the [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/). Again, this will also be validated by our CI, but checking it yourself beforehand will speed things up drastically.

## Maintainer Team at Merico

A group of engineers maintains DevStream at Merico, led by [@ironcore864](https://github.com/ironcore864).

We aim to reply to issues within 24 hours.
