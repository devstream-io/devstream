# [[.AppName]]

This is a repo for app [[.AppName]]; bootstrapped by DevStream.

By default, the automatically generated scaffolding contains:

- a piece of sample go cli app code using the [Cobra Commander Framework](https://github.com/spf13/cobra)
- [GoReleaser](https://goreleaser.com/) for building and releasing the app
- .gitignore
- Makefile

## Automatic Releases

Just push a tag to the repo and [GoReleaser](https://goreleaser.com/) will build and release the app.

For example, to release version 0.1.0, run:

```git
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```
