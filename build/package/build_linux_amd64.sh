#!/usr/bin/env bash

set -e

export GOOS=linux
export GOARCH=amd64
export VERSION
go build -trimpath -gcflags="all=-N -l" -o output/dtm-${GOOS}-${GOARCH} ./cmd/devstream/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/golang
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/python
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/nodejs
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trellogithub/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocd-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocd/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocdapp/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/jenkins-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jenkins/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/kubeprometheus/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/reposcaffolding/github/golang/
go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/devlake-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/devlake/
