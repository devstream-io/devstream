<div align="center">
<br/>
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

<br/>

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![Test](https://github.com/merico-dev/stream/actions/workflows/main-builder.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/stream)](https://goreportcard.com/report/github.com/merico-dev/stream)

# DevStream
</div>

## DevStream, What Is It Anyway?

DevStream is an open-source DevOps toolchain manager.

Imagine you are in a new project. Before writing the first line of code, you would have to figure out the tools needed in the whole Software Development Life Cycle (SDLC). You would probably need the following pieces:

- some kind of project management software or issue tracking tools (e.g., Jira);
- 
- someplace for source code management (GitHub and alike);
- 
- some tools for continuous integration (e.g., Jenkins, GitHub Actions, CircleCI, Travis CI);
- some tools for continuous delivery or continuous deployment (e.g., fluxcd/flux2, ArgoCD, etc.);
- someplace serving as the single source of truth for secrets and credentials (secrets manager, e.g., Vault by HashiCorp);
- some tools for centralized logging and monitoring (for example, ELK, Prometheus/Grafana);

And maybe more. The list could go on for quite a bit.

And, there are multiple challenges in creating YOUR ideal SDLC workflow:

- There are too many choices. Even for a particular field, there are too many. Which is best? There is no "one-size-fits-all" answer because it totally depends.
- Integration between different pieces.
- The software world (and the DevOps world) changes and it changes fast. What's best for today might not be the best tomorrow. You want to switch some parts out and get some new state-of-the-art pieces in so that you always keep your efficiency high.

To be fair, there are a few integrated products out there that may contain everything you might need, but they might not suit your specific requirements perfectly. So, the chance is, you will still want to go out and do your research, find the best pieces for you, and integrate them. And, it would be a lot of operational overhead if all you had to do all day was install and uninstall and integrate things.

You probably have already seen where we are going with this, and you are right: DevStream, an open-source DevOps toolchain manager, aims to be the solution here.

Think of the Linux kernel V.S. different distributions. Different distros offer different packages so that you can always choose the best for your need.

Or, Think of `yum`, `apt`, or `apk`. You can easily set it up with your favorite packages for any new environment using these package managers.

DevStream aims to be the package manager for DevOps tools. To be more ambitious, DevStream wants to be the Linux kernel, around which different distros can be created with various components so that you can always have the best components for each part of your SDLC workflow.

## Why Using DevStream?

No more manual curl/wget download, apt install, helm install; no more local experiments and playing around just to get a piece of tool installed correctly.

Define your wanted DevOps tools in a single human-readable YAML config file, and at the press of a button (one single command), you will have your whole DevOps toolchain and SDLC workflow set up.

Want to install another different tool for a try? No problem. Want to remove or reinstall a specific piece in the whole workflow? Got your back.

## Architecture

See [docs/architecture.md](./docs/architecture.md).

## Build

```shell
$ make help

Usage:
  make <target>
  help                Display this help.
  build               Build dtm & plugins locally.
  build-core          Build dtm core only, without plugins, locally.
  build-release       Build for all platforms for release.
  build-darwin-arm64  Build for darwin/arm64 for release.
  build-darwin-amd64  Cross-platform build for darwin/amd64.
  build-linux-amd64   Cross-platform build for linux/amd64
  e2e                 Run e2e tests.
  fmt                 Run 'go fmt' & goimports against code.
  vet                 Run go vet against code.
```

## Test

Unit/functional test:

```bash
go test ./...
```

## Configuration

See [examples/config.yaml](./examples/config.yaml).

## Run

To install/reinstall/update, run:

```bash
./dtm apply -f examples/config.yaml
```

To delete, run:

```bash
./dtm delete -f examples/config.yaml
```

To verify, run:

```bash
./dtm verify -f examples/config.yaml
```

## Why `dtm`?

Q: The CLI tool is named `dtm`, while the tool itself is called DevStream. What the heck?! Where is the consistency?

A: Inspired by [`git`](https://github.com/git/git#readme), the name is (depending on your mood):

- a symmetric, scientific acronym of **d**evs**t**rea**m**.
- "devops toolchain manager": you're in a good mood, and it actually works for you.
- "dead to me": when it breaks.

## Contribute

See [CONTRIBUTING.md](./CONTRIBUTING.md).
