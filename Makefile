MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH=$(patsubst %/,%,$(dir $(MKFILE_PATH)))/build/working_dir
VERSION=0.2.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: fmt vet ## Build dtm & plugins locally.
	go mod tidy
	mkdir -p .devstream
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/golang
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-python-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/python
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-nodejs-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/githubactions/nodejs
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-github-integ-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/trellogithub/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/argocdapp/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/jenkins-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/jenkins/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/kube-prometheus-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/kubeprometheus/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/github-repo-scaffolding-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/reposcaffolding/github/golang/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/devlake-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/devlake/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/gitlabci-golang-${GOOS}-${GOARCH}_${VERSION}.so ./cmd/gitlabci/golang
	go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/

build-core: fmt vet ## Build dtm core only, without plugins, locally.
	go mod tidy
	go build -trimpath -gcflags="all=-N -l"  -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm ./cmd/devstream/

clean:
	rm -rf .devstream
	rm -f dtm*
	rm -rf build/working_dir

build-linux-amd64: ## Cross-platform build for linux/amd64
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
	goimports -local="github.com/merico-dev/stream" -d -w internal
	goimports -local="github.com/merico-dev/stream" -d -w test
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

e2e: build ## Run e2e tests.
	./dtm apply -f config.yaml
	./dtm verify -f config.yaml
	./dtm delete -f config.yaml

e2e-up: ## Start kind cluster for e2e tests
	sh hack/e2e/e2e-up.sh

e2e-down: ## Stop kind cluster for e2e tests
	sh hack/e2e/e2e-down.sh
