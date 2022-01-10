# DevStream hack GuideLines

This document describes how you can use the scripts from [`hack`](.) directory
and gives a brief introduction and explanation of these scripts.

## Overview

The [`hack`](.) directory contains many scripts that ensure continuous development of DevStream.

## Key scripts

- [e2e](./e2e): This directory holds the scripts used for e2e testing.
  - [e2e-up.sh](./e2e/e2e-up.sh): This script used for setup e2e testing environment.
  - [e2e-down.sh](./e2e/e2e-down.sh): This script used for clear e2e testing kind cluster.

## Examples

- Setup e2e testing environment

```shell
cd hack/e2e
sh e2e-up.sh
```
