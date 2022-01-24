#!/usr/bin/env bash

set -e

export GOOS=linux
export GOARCH=amd64
VERSION
go build -trimpath -gcflags="all=-N -l" -o output/dtm-${GOOS}-${GOARCH} ./cmd/devstream/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trellogithub/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocd-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocd/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocdapp/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/kube-prometheus/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/reposcaffolding/github/golang/
