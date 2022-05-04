
SELF_DIR=$(dir $(lastword $(MAKEFILE_LIST)))

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GO_PLUGIN_BUILD=go build -buildmode=plugin -trimpath -gcflags="all=-N -l"
PLUGINS=$(notdir $(wildcard $(ROOT_DIR)/cmd/plugin/*))
PLUGIN_SUFFIX=${GOOS}-${GOARCH}_${VERSION}

DTM_ROOT=github.com/devstream-io/devstream
GO_LDFLAGS += -X '$(DTM_ROOT)/internal/pkg/version.Version=$(VERSION)' \
		-X '$(DTM_ROOT)/cmd/devstream/list.PluginsName=$(PLUGINS)'

FIND := find . -path './cmd/**/*.go' -o -path './test/**/*.go' -o -path './pkg/**/*.go' -o -path './internal/**/*.go'

ifeq ($(GOOS),linux)
	MD5SUM=md5sum
else
	MD5SUM=md5 -q
endif

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(SELF_DIR) && pwd -P))
endif

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

ifeq ($(origin PLUGINS_DIR),undefined)
PLUGINS_DIR := $(ROOT_DIR)/.devstream
$(shell mkdir -p $(PLUGINS_DIR))
endif

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9.%-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Remove dtm and plugins. It's best to run a "clean" before "build".
	-rm -rf .devstream
	-rm -f dtm*

.PHONY: build-core
build-core: fmt vet mod-tidy ## Build dtm core only, without plugins, locally.
	go build -trimpath -gcflags="all=-N -l" -ldflags "$(GO_LDFLAGS)" -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/
	@-rm -f dtm
	@mv dtm-${GOOS}-${GOARCH} dtm
	@echo ">>>>>>>>>>>> 'dtm' has been generated in the current directory"

.PHONY: build-plugin.%
build-plugin.%: fmt vet mod-tidy ## Build one dtm plugin, like "make build-plugin.argocd".
	$(eval plugin_name := $(strip $*))
	${GO_PLUGIN_BUILD} -o .devstream/${plugin_name}-${PLUGIN_SUFFIX}.so ${ROOT_DIR}/cmd/plugin/${plugin_name}
	@$(MAKE) md5-plugin.$(plugin_name)

.PHONY: build-plugins
build-plugins: fmt vet mod-tidy $(addprefix build-plugin.,$(PLUGINS)) ## Build dtm plugins only. Use multi-threaded like "make build-plugins -j8" to speed up.

.PHONY: build
build: build-core build-plugins ## Build everything. Use multi-threaded like "make build -j8" to speed up.

.PHONY: md5
md5: md5-plugins ## Create md5 sums for all plugins.

.PHONY: md5-plugins
md5-plugins: $(addprefix md5-plugin.,$(PLUGINS))

.PHONY: md5-plugin.%
md5-plugin.%:
	$(eval plugin_name := $(strip $*))
	${MD5SUM} .devstream/${plugin_name}-${PLUGIN_SUFFIX}.so > .devstream/${plugin_name}-${PLUGIN_SUFFIX}.md5

.PHONY: fmt
fmt:  ## Run 'go fmt' & goimports against code.
	@echo ">>>>>>>>>>>> Formating codes"
	@go install golang.org/x/tools/cmd/goimports@latest
	@$(FIND) -type f | xargs gofmt -s -w
	@$(FIND) -type f | xargs goimports -w -local $(DTM_ROOT)

.PHONY: vet
vet: ## Run "go vet ./...".
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Run "go mod tidy".
	go mod tidy

.PHONY: e2e
e2e: build ## Run e2e tests.
	./dtm apply -f config.yaml
	./dtm verify -f config.yaml
	./dtm delete -f config.yaml

.PHONY: e2e-up
e2e-up: ## Start kind cluster for e2e tests.
	sh hack/e2e/e2e-up.sh

.PHONY: e2e-down
e2e-down: ## Stop kind cluster for e2e tests.
	sh hack/e2e/e2e-down.sh
