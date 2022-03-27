#!/usr/bin/env bash

set -e

export GOOS=linux
export GOARCH=amd64
export VERSION
export PLUGINS_CMD_ROOT=./cmd/plugin
export PLUGIN_SUFFIX=${GOOS}-${GOARCH}_${VERSION}
# buid core
go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o output/dtm-${GOOS}-${GOARCH} ./cmd/devstream/

# buid plugins
PLUGINS_DIR=$(find ${PLUGINS_CMD_ROOT} -name "main.go" -exec dirname {} \;)
for plugin_dir in ${PLUGINS_DIR}; do 
    plugin_name="${plugin_dir##*/}"
    go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/"${plugin_name}"-"${PLUGIN_SUFFIX}".so ${PLUGINS_CMD_ROOT}/"${plugin_name}"
	md5sum output/"${plugin_name}"-"${PLUGIN_SUFFIX}".so > output/"${plugin_name}"-"${PLUGIN_SUFFIX}".md5
done
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions-golang
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions-python
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions-nodejs
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/trello-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trello
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trellogithub
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/jira-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jiragithub
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocd-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocd
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocdapp
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/jenkins-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jenkins
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/kubeprometheus
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/github-repo-scaffolding-golang
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/devlake-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/devlake
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/gitlabci-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/gitlabci/golang
# go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o output/openldap-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/openldap
