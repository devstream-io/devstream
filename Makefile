MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH=$(patsubst %/,%,$(dir $(MKFILE_PATH)))/build/working_dir
VERSION=0.3.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

ifeq ($(GOOS),linux)
  MD5SUM=md5sum
else
  MD5SUM=md5 -q
endif

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: build-core build-plugins md5 ## Build dtm core & plugins locally.
	rm -f dtm
	mv dtm-${GOOS}-${GOARCH} dtm

build-core: fmt vet ## Build dtm core only, without plugins, locally.
	go mod tidy
	go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/cmd/devstream/version.Version=${VERSION}" -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/

build-plugins: fmt vet ## Build dtm plugins only, without core, locally.
	go mod tidy
	mkdir -p .devstream
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/golang
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/python
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/nodejs
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trello/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trellogithub/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocdapp/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/jenkins-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jenkins/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/kubeprometheus/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/reposcaffolding/github/golang/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/devlake-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/devlake/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/gitlabci-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/gitlabci/golang
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/jira-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jiragithub/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/openldap-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/openldap/

md5: md5-core md5-plugins

md5-core:
	${MD5SUM} dtm-${GOOS}-${GOARCH} > dtm-${GOOS}-${GOARCH}.md5

md5-plugins:
	${MD5SUM} .devstream/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/trello-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/trello-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/argocd-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/argocd-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/argocdapp-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/jenkins-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/jenkins-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/devlake-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/devlake-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/gitlabci-golang-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/gitlabci-golang-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/jira-github-integ-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/jira-github-integ-${GOOS}-${GOARCH}_${VERSION}.md5
	${MD5SUM} .devstream/openldap-${GOOS}-${GOARCH}_${VERSION}.so > .devstream/openldap-${GOOS}-${GOARCH}_${VERSION}.md5

clean: ## Remove local plugins and locally built artifacts.
	rm -rf .devstream
	rm -f dtm*
	rm -rf build/working_dir

build-linux-amd64: ## Cross-platform build for "linux/amd64".
	echo "Building in ${BUILD_PATH}"
	mkdir -p .devstream
	rm -rf ${BUILD_PATH} && mkdir ${BUILD_PATH}
	docker buildx build --platform linux/amd64 --load -t mericodev/stream-builder:v${VERSION} --build-arg VERSION=${VERSION} -f build/package/Dockerfile .
	cp -r go.mod go.sum cmd internal pkg build/package/build_linux_amd64.sh ${BUILD_PATH}/
	chmod +x ${BUILD_PATH}/build_linux_amd64.sh
	docker run --rm --platform linux/amd64 -v ${BUILD_PATH}:/devstream mericodev/stream-builder:v${VERSION}
	mv ${BUILD_PATH}/output/*.so .devstream/
	mv ${BUILD_PATH}/output/dtm* .
	rm -rf ${BUILD_PATH}

fmt: ## Run 'go fmt' & goimports against code.
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -local="github.com/merico-dev/stream" -d -w cmd
	goimports -local="github.com/merico-dev/stream" -d -w pkg
	goimports -local="github.com/merico-dev/stream" -d -w internal
	goimports -local="github.com/merico-dev/stream" -d -w test
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

e2e: build ## Run e2e tests.
	./dtm-${GOOS}-${GOARCH} apply -f config.yaml
	./dtm-${GOOS}-${GOARCH} verify -f config.yaml
	./dtm-${GOOS}-${GOARCH} delete -f config.yaml

e2e-up: ## Start kind cluster for e2e tests.
	sh hack/e2e/e2e-up.sh

e2e-down: ## Stop kind cluster for e2e tests.
	sh hack/e2e/e2e-down.sh
