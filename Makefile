MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH=$(patsubst %/,%,$(dir $(MKFILE_PATH)))/build/working_dir
VERSION=0.0.1

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: fmt vet ## Build dtm & plugins locally.
	go mod tidy
	mkdir -p .devstream
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions_${VERSION}.so ./cmd/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-github-integ_${VERSION}.so ./cmd/trellogithub/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd_${VERSION}.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp_${VERSION}.so ./cmd/argocdapp/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/kube-prometheus_${VERSION}.so ./cmd/kubeprometheus/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/github-repo-scaffolding-golang_${VERSION}.so ./cmd/reposcaffolding/github/golang/
	go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm ./cmd/devstream/

build-core: fmt vet ## Build dtm core only, without plugins, locally.
	go mod tidy
	go build -trimpath -gcflags="all=-N -l"  -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm ./cmd/devstream/

clean:
	rm -rf .devstream
	rm -f dtm*
	rm -rf build/working_dir

build-release: build-darwin-arm64 build-darwin-amd64 build-linux-amd64 ## Build for all platforms for release.

build-darwin-arm64: ## Build for darwin/arm64 for release.
	go mod tidy
	mkdir -p .devstream
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-darwin-arm64_${VERSION}.so ./cmd/githubactions/
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-github-integ-darwin-arm64_${VERSION}.so ./cmd/trellogithub/
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd-darwin-arm64_${VERSION}.so ./cmd/argocd/
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp-darwin-arm64_${VERSION}.so ./cmd/argocdapp/
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/kube-prometheus-darwin-arm64_${VERSION}.so ./cmd/kubeprometheus/
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/github-repo-scaffolding-golang-darwin-arm64_${VERSION}.so ./cmd/reposcaffolding/github/golang/
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm-darwin-arm64 ./cmd/devstream/

build-darwin-amd64: ## Cross-platform build for darwin/amd64.
	go mod tidy
	mkdir -p .devstream
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions-darwin-amd64_${VERSION}.so ./cmd/githubactions/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/trello-github-integ-darwin-amd64_${VERSION}.so ./cmd/trellogithub/
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd-darwin-amd64_${VERSION}.so ./cmd/argocd/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp-darwin-amd64_${VERSION}.so ./cmd/argocdapp/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/kube-prometheus-darwin-amd64_${VERSION}.so ./cmd/kubeprometheus/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/github-repo-scaffolding-golang-darwin-amd64_${VERSION}.so ./cmd/reposcaffolding/github/golang/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.VERSION=${VERSION}" -o dtm-darwin-amd64 ./cmd/devstream/

build-linux-amd64: ## Cross-platform build for linux/amd64
	echo "Building in ${BUILD_PATH}"
	mkdir -p .devstream
	rm -rf ${BUILD_PATH} && mkdir ${BUILD_PATH}
	docker buildx build --platform linux/amd64 --load -t mericodev/stream-builder:v${VERSION} --build-arg http_proxy=${VERSION} -f build/package/Dockerfile .
	cp -r go.mod go.sum cmd internal build/package/build_linux_amd64.sh ${BUILD_PATH}/
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
