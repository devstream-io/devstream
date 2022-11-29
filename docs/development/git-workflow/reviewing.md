# Code Reviewing Guide

This document covers who may review pull requests for this project, and provides guidance on how to perform code reviews that meet our community standards and code of conduct. All reviewers must read this document and agree to follow the project review guidelines. Reviewers who do not follow these guidelines may have their privileges revoked.

[howto]: https://contribute.cncf.io/maintainers/github/templates/recommended/reviewing/

## The Reviewer Role

The reviewer role is distinct from the maintainer role.

Reviewers can approve a pull request but they cannot merge it. A maintainer handles final approval and merging the pull request.

## Values

All reviewers must abide by the [Code of Conduct](https://github.com/devstream-io/devstream/blob/main/CODE_OF_CONDUCT.md) and are also protected by it. A reviewer should not tolerate poor behavior and is encouraged to report any behavior that violates the Code of Conduct. All of our values listed above are distilled from our Code of Conduct.

Below are concrete examples of how it applies to code review specifically:

### Inclusion

Be welcoming and inclusive. You should proactively ensure that the author is successful. While any particular pull request may not ultimately be merged, overall we want people to have a great experience and be willing to contribute again. Answer the questions they didn't know to ask or offer concrete help when they appear stuck.

### Sustainability

Avoid burnout by enforcing healthy boundaries. Here are some examples of how a reviewer is encouraged to act to take care of themselves:

- Authors should meet baseline expectations when submitting a pull request, such as writing tests.
- If your availability changes, you can step down from a pull request and have someone else assigned.
- If interactions with an author are not following code of conduct, close the PR and raise it up with your Code of Conduct committee or point of contact. It's not your job to coax people into behaving.

### Trust

Be trustworthy. During a review, your actions both build and help maintain the trust that the community has placed in this project. Below are examples of ways that we build trust:

- **Transparency** - If a pull request won't be merged, clearly say why and close it. If a pull request won't be reviewed for a while, let the author know so they can set expectations and understand why it's blocked.
- **Integrity** - Put the project's best interests ahead of personal relationships or company affiliations when deciding if a change should be merged.
- **Stability** - Only merge when then change won't negatively impact project stability. It can be tempting to merge a pull request that doesn't meet our quality standards, for example when the review has been delayed, or because we are trying to deliver new features quickly, but regressions can significantly hurt trust in our project.

## Process

- A PR might be automatically assigned according to the OWNERS file. If not, any reviewer can assign the pull request, and set specific labels.
- The project uses GitHub actions for automation, but currently, no extra automation/bots are used for assigning PRs to reviewers.
- The reviewer should wait for automated checks to pass before reviewing.
- Docs/automation related critical changes require reviews from maintainers. Minor tweaks which are clear enough don't require maintainers' review. When in doubt, ask the maintainers to double-check.
- GitHub Actions checks must pass before merging, with the only exception being the commit message lint, which can be improved by squashing and re-editing the message.
- When a PR is stuck, reviewers should either reach out to the contributor to help to move things forward, or report the status to the maintainers.
- In general, it's not recommended for reviewers to commit changes directly to the pull request. Exceptions are: when the original author has abandoned the pull request; the original author is a new contributor and might need help.
- Maintainers can merge their own pull requests after it has been reviewed.
- Maintainers can merge pull requests without review in times of great need.

## Checklist

Below are a set of common questions that apply to all pull requests:

- [ ] Is this PR targeting the correct branch?
- [ ] Does the commit message provide an adequate description of the change?
- [ ] Does the commit message follow the [conventional commit messages specification](https://www.conventionalcommits.org/en/v1.0.0/)?
- [ ] Does the affected code have corresponding tests?
- [ ] Are the changes documented, not just with inline documentation, but also with conceptual documentation such as an overview of a new feature, or task-based documentation like a tutorial? Consider if this change should be announced on DevStream's blog.
- [ ] Does this introduce breaking changes that would require an announcement or bumping the major version?

## Reading List

Reviewers are encouraged to read the following articles for help with common reviewer tasks:

- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [Project Layout](https://docs.devstream.io/en/latest/development/project-layout/)
- [The Art of Closing: How to closing an unfinished or rejected pull request](https://blog.jessfraz.com/post/the-art-of-closing/)
- [Kindness and Code Reviews: Improving the Way We Give Feedback](https://product.voxmedia.com/2018/8/21/17549400/kindness-and-code-reviews-improving-the-way-we-give-feedback)
- [Code Review Guidelines for Humans: Examples of good and back feedback](https://phauer.com/2018/code-review-guidelines/#code-reviews-guidelines-for-the-reviewer)
