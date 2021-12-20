build: ## Build dtm & plugins locally.
	go get ./...
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/githubactions_0.0.1.so ./cmd/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocd_0.0.1.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocdapp_0.0.1.so ./cmd/argocdapp/
	go build -trimpath -gcflags="all=-N -l" -o dtm ./cmd/devstream/

core: ## Build dtm core only, locally.
	go get ./...
	go build -trimpath -gcflags="all=-N -l" -o dtm ./cmd/devstream/

fmt: ## Run go fmt against code.
	goimports -local="github.com/merico-dev/stream" -d -w cmd
	goimports -local="github.com/merico-dev/stream" -d -w internal
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...
