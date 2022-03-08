<div align="center">
<br/>
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

<br/>

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![Test](https://github.com/merico-dev/stream/actions/workflows/main-builder.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/stream)](https://goreportcard.com/report/github.com/merico-dev/stream)
[![Downloads](https://img.shields.io/github/downloads/merico-dev/stream/total.svg)](https://github.com/merico-dev/stream/releases)
  
# DevStream
</div>

## DevStream, What Is It Anyway?

> English | [中文](docs/README_zh.md)

TL;DR: DevStream (CLI tool named `dtm`) is an open-source DevOps toolchain manager.

Imagine you are starting a new project or ramping up a new team. Before writing the first line of code, you have to figure out the tools to run an effective SDLC process and from development to deployment. 

Typically, you'd need the following pieces in place to work effectively:

- Project management software or issue tracking tools (JIRA, etc.)
- Source code management (GitHub, Bitbucket, etc.)
- Continuous integration tools (Jenkins, CircleCI, Travis CI, etc.)
- Continuous delivery/deployment tools (fluxcd/flux2, ArgoCD, etc.)
- A single source of truth for secrets and credentials (secrets manager, e.g., Vault by HashiCorp)
- Some tools for centralized logging and monitoring (for example, ELK, Prometheus/Grafana);

The list could go on for quite a bit, but you get the idea!

There are many challenges in creating an effective and personalized workflow:

- There are too many choices. Which is best? There is no "one-size-fits-all" answer because it totally depends on your needs and preferences.
- Integration between different pieces is challenging, creating silos and fragmentation.
- The software world evolves fast. What's best today might not make sense tomorrow. If you want to switch parts or tools out, it can be challenging and resource intensive to manage.

To be fair, there are a few integrated products out there that may contain everything you might need, but they might not suit your specific requirements perfectly. So, the chances are, you will still want to go out and do your research, find the best pieces, and integrate them yourself. That being said, to choose, launch, connect, and manage all these pieces take a lot of time and energy.

You might be seeing where we are going with this...

We wanted to make it easy to set up these personalized and flexible toolchains, so we built DevStream, an open-source DevOps toolchain manager.

Think of the Linux kernel V.S. different distributions. Different distros offer different packages so that you can always choose the best for your need.

Or, think of `yum`, `apt`, or `apk`. You can easily set it up with your favorite packages for any new environment using these package managers.

DevStream aims to be the package manager for DevOps tools.

To be more ambitious, DevStream wants to be the Linux kernel, around which different distros can be created with various components so that you can always have the best components for each part of your SDLC workflow.

## Why Use DevStream?

No more manual curl/wget download, apt install, helm install; no more local experiments and playing around just to get a piece of tool installed correctly.

Define your desired DevOps tools in a single human-readable YAML config file, and at the press of a button (one single command), you will have your whole DevOps toolchain and SDLC workflow set up. Five Minutes. One Command.

Want to install another different tool for a try? No problem.

Want to remove or reinstall a specific piece in the workflow? DevStream has got your back!

## Supported DevOps Tools

| Type                   | Plugin                         | Note                           | Usage/Doc |
|------------------------|--------------------------------|--------------------------------|-----------|
| Issue Tracking         | trello-github-integ            | Trello/GitHub integration      | [doc](./docs/plugins/trello-github-integ_plugin.md) |
| Source Code Management | github-repo-scaffolding-golang | Go WebApp scaffolding          | [doc](./docs/plugins/github-repo-scaffolding-golang_plugin.md) |
| CI                     | jenkins                        | Jenkins installation           | [doc](./docs/plugins/jenkins_plugin.md) |
| CI                     | githubactions-golang           | GitHub Actions CI for Golang   | [doc](./docs/plugins/githubactions-golang_plugin.md)          |
| CI                     | githubactions-python           | GitHub Actions CI for Python   | [doc](./docs/plugins/githubactions-python_plugin.md)          |
| CI                     | githubactions-nodejs           | GitHub Actions CI for Nodejs   | [doc](./docs/plugins/githubactions-nodejs_plugin.md)          |
| CI                     | gitlabci-golang                | GitLab CI for Golang           | [doc](./docs/plugins/gitlabci-golang_plugin.md)          |
| CD/GitOps              | argocd                         | ArgoCD installation            | [doc](./docs/plugins/argocd_plugin.md)          |
| CD/GitOps              | argocdapp                      | ArgoCD Application creation    | [doc](./docs/plugins/argocdapp_plugin.md)          |
| Monitoring             | kube-prometheus                | Prometheus/Grafana K8s install | [doc](./docs/plugins/kube-prometheus_plugin.md)          |
| Observability          | devlake                        | DevLake installation           | [doc](./docs/plugins/devlake_plugin.md)          |
| LDAP                   | openldap                       | OpenLDAP installation          | [doc](./docs/plugins/openldap_plugin.md)          |

## Quick Start

If you want to get a quick start, follow our [quick start](./docs/quickstart_en.md) doc now.

## Configuration

This is an example of DevStream config: [examples/quickstart.yaml](./examples/quickstart.yaml).

Remember to open this configuration file, modify all FULL_UPPER_CASE_STRINGS (like YOUR_GITHUB_USERNAME, for example) in it to your own.

Pay attention to the meaning of each item to ensure that it is what you want.

For other plugins, checkout the [docs/plugins](./docs/plugins/) folder for detailed usage in their documentations.

## Run

To apply the config, run:

```bash
./dtm apply -f YOUR_CONFIG_FILE.yaml
```

If you don't specify the config file with the "-f" parameter, it will try to use the default value which is "config.yaml" from the current directory.

_`dtm` will compare the config, the state, and the resources to decide whether a "create", "update", or "delete" is needed. For more information, see [this document here](./docs/config_state_resource_explanation.md)._

The command above will ask you for confirmation before actually executing the changes. To apply without confirmation (like `apt-get -y update`), run:

```bash
./dtm -y apply -f YOUR_CONFIG_FILE.yaml
```

To delete everything defined in the config, run:

```bash
./dtm delete -f YOUR_CONFIG_FILE.yaml
```

_Note that this deletes everything defined in the config. If some config is deleted after apply (state has it but config not), `dtm delete` won't delete it. It differs from `dtm destroy`._

Similarly, to delete without confirmation:

```bash
./dtm -y delete -f YOUR_CONFIG_FILE.yaml
```
To delete everything defined in the config, regardless of the state: 

```bash
./dtm delete --force -f YOUR_CONFIG_FILE.yaml
```

To verify, run:

```bash
./dtm verify -f YOUR_CONFIG_FILE.yaml
```

To destroy everything, run:

```bash
./dtm destroy
```

_`dtm` will read the state, then determine which tools are installed, and then remove those tools. It's same as `dtm apply -f empty.yaml` (empty.yaml is an empty config file)._

## Dev Info

### Source

#### Prerequisite Tools

- Git
- Go (1.17+)

#### Fetch from GitHub

```bash
mkdir -p ~/gocode
cd ~/gocode
git clone https://github.com/merico-dev/stream.git
```

#### Build

```bash
cd ~/gocode/stream
make build
mv dtm-$(go env GOOS)-$(go env GOARCH) dtm
```

See the Makefile for more info.

```makefile
$ make help

Usage:
  make <target>
  help                Display this help.
  build               Build dtm & plugins locally.
  build-core          Build dtm core only, without plugins, locally.
  build-linux-amd64   Cross-platform build for linux/amd64
  fmt                 Run 'go fmt' & goimports against code.
  vet                 Run go vet against code.
  e2e                 Run e2e tests.
  e2e-up              Start kind cluster for e2e tests
  e2e-down            Stop kind cluster for e2e tests
```

#### Test

Run unit tests:

```bash
go test ./...
```

Run e2e tests:

```bash
make e2e
```

## Architecture

See [docs/architecture.md](./docs/architecture.md).

## Why `dtm`?

Q: The CLI tool is named `dtm`, while the tool itself is called DevStream. What the heck?! Where is the consistency?

A: Inspired by [`git`](https://github.com/git/git#readme), the name is (depending on your mood):

- a symmetric, scientific acronym of **d**evs**t**rea**m**.
- "devops toolchain manager": you're in a good mood, and it actually works for you.
- "dead to me": when it breaks.

## Community

- <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>: Message us on Discord

## Contribute

See [CONTRIBUTING.md](./CONTRIBUTING.md).
