# Contributing to DevStream

Thanks for contributing, and thanks for reading this doc before doing it!

This is a set of guidelines for contributing to DevStream. These are only guidelines, not rules. So, use your best judgment, and feel free to propose changes to this document.

## How Can I Contribute?

Report bugs, create issues, propose new features by creating an issue with the corresponding issue template, and label it accordingly.

If you intend to change the public API or make any non-trivial change, filing an issue first is recommended. This lets us reach an agreement on your proposal before putting significant effort into it.

If you are only fixing a bug, it’s OK to submit a PR without an issue. In this case, Filing an issue is still recommended because it helps us track the issue.

## Maintainer Team at Merico

A group of engineers maintains DevStream at Merico, led by [@ironcore864](https://github.com/ironcore864).

We aim to reply to issues within 24 hours.

## Contributor Ladder Growth Programs

See [contributor_ladder_growth_programs.md](docs/contributing/contributor_ladder_growth_programs.md).

## Style Guides

### Linters

We use `golangci-lint` ([official website](https://golangci-lint.run/), [GitHub](https://github.com/golangci/golangci-lint)) for linting, which is a Go linters aggregator. It's also integrated with the GitHub Actions workflows.

- The list of linters enabled by default is documented [here](https://golangci-lint.run/usage/linters/).
- It can be [integrated with IDE](https://golangci-lint.run/usage/integrations/)
- You can [run it locally](https://golangci-lint.run/usage/quick-start/).

Besides, we also use the [Go Report Card](https://goreportcard.com/report/github.com/devstream-io/devstream). There is a badge like [![Go Report Card](https://goreportcard.com/badge/github.com/devstream-io/devstream)](https://goreportcard.com/report/github.com/devstream-io/devstream) in the main README.md.

### Git Commit Message

We try our best to follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) guidelines.

TL;DR: The commit message should be structured as follows:

```
<type>[optional scope]: <description>
[optional body]
[optional footer(s)]
```

where "type" can be:
- `feat`: implements a feature
- `fix`: patches a bug
- `BREAKING CHANGE`: a breaking change. Or append `!` at the end of "feat" or "fix", like `feat!` and `fix!`
- other types are allowed, for example: `build`, `chore`, `ci`, `docs`, `style`, `refactor`, `perf`, `test`

Both "body" and "footer" are optional; BREAKING CHANGE can be addressed both in the title as well as in the footer. Some examples:

- `feat: send an email to the customer when a product is shipped`
- `feat!: send an email to the customer when a product is shipped`
- `feat(api): send an email to the customer when a product is shipped`
- `feat(api)!: send an email to the customer when a product is shipped`
- `BREAKING CHANGE: send an email to the customer when a product is shipped`
- ```
  feat!: send an email to the customer when a product is shipped
  A detailed description in the body.
  BREAKING CHANGE: readdressing the breaking change in the footer.
  ```

### Creating a New Plugin

See our [doc about creating a plugin here](https://www.devstream.io/docs/creating-a-plugin).

## Development Workflow

See [development workflow documentation here](https://www.devstream.io/docs/development-workflow).
