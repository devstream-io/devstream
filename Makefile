help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: fmt vet ## Build dtm & plugins locally.
	go mod tidy
	mkdir -p .devstream
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/githubactions_0.0.1.so ./cmd/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocd_0.0.1.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o .devstream/argocdapp_0.0.1.so ./cmd/argocdapp/
	go build -trimpath -gcflags="all=-N -l" -o dtm ./cmd/devstream/

build-core: fmt vet ## Build dtm core only, locally.
	go mod tidy
	go build -trimpath -gcflags="all=-N -l" -o dtm ./cmd/devstream/

fmt: ## Run 'go fmt' & goimports against code.
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -local="github.com/merico-dev/stream" -d -w cmd
	goimports -local="github.com/merico-dev/stream" -d -w internal
	goimports -local="github.com/merico-dev/stream" -d -w test
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

test-cmd: build ## Run e2e tests against the cmd.
	go build -trimpath -gcflags="all=-N -l" -o smoke ./test/smoke/
	./smoke
