#!/usr/bin/env bash

set -e

export GOOS=linux
export GOARCH=amd64
go build -trimpath -gcflags="all=-N -l" -o output/dtm-${GOOS}-${GOARCH} ./cmd/devstream/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-${GOOS}-${GOARCH}_0.0.1.so ./cmd/githubactions/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/trello-github-integ-${GOOS}-${GOARCH}_0.0.1.so ./cmd/trellogithub/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocd-${GOOS}-${GOARCH}_0.0.1.so ./cmd/argocd/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocdapp-${GOOS}-${GOARCH}_0.0.1.so ./cmd/argocdapp/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/kube-prometheus-${GOOS}-${GOARCH}_0.0.1.so ./cmd/kube-prometheus/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/github-repo-scaffolding-${GOOS}-${GOARCH}_0.0.1.so ./cmd/reposcaffolding/github/
