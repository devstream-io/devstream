#!/usr/bin/env bash

set -e

GOOS=${GOOS:-linux} GOARCH=${GOARCH:-amd64} go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/plugins/githubactions_0.0.1.so ./cmd/githubactions/
GOOS=${GOOS:-linux} GOARCH=${GOARCH:-amd64} go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/plugins/argocd_0.0.1.so ./cmd/argocd/
GOOS=${GOOS:-linux} GOARCH=${GOARCH:-amd64} go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/plugins/argocdapp_0.0.1.so ./cmd/argocdapp/
GOOS=${GOOS:-linux} GOARCH=${GOARCH:-amd64} go build -trimpath -gcflags="all=-N -l" -o output/dsm ./cmd/devstream/
