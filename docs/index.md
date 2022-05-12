# DevStream Docs

## DevStream, What Is It Anyway?

TL;DR: DevStream (CLI tool named `dtm`) is an open-source DevOps toolchain manager.

Imagine you are starting a new project or ramping up a new team. Before writing the first line of code, you have to figure out the tools to run an effective SDLC process and from development to deployment.

Typically, you'd need the following pieces in place to work effectively:

- Project management software or issue tracking tools (JIRA, etc.)
- Source code management (GitHub, Bitbucket, etc.)
- Continuous integration tools (Jenkins, CircleCI, Travis CI, etc.)
- Continuous delivery/deployment tools (Flux CD/Flux2, Argo CD, etc.)
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

## Why `dtm`?

Q: The CLI tool is named `dtm`, while the tool itself is called DevStream. What the heck?! Where is the consistency?

A: Inspired by [`git`](https://github.com/git/git#readme), the name is (depending on your mood):

- a symmetric, scientific acronym of **d**evs**t**rea**m**.
- "devops toolchain manager": you're in a good mood, and it actually works for you.
- "dead to me": when it breaks.
