# Commit Messages

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
