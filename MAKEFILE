SELF_DIR=$(dir $(lastword $(MAKEFILE_LIST)))
GOOS=$(shell go env GOOS)
GOPATH=$(shell go env GOPATH)
GOARCH=$(shell go env GOARCH)
GO_PLUGIN_BUILD=go build -buildmode=plugin -trimpath -gcflags="all=-N -l"
PLUGINS=$(notdir $(wildcard $(ROOT_DIR)/cmd/plugin/*))
PLUGIN_SUFFIX=${GOOS}-${GOARCH}_${VERSION}
SHELL := /bin/bash
DTM_ROOT=github.com/devstream-io/devstream
GO_LDFLAGS += -X '$(DTM_ROOT)/internal/pkg/version.Version=$(VERSION)' \
		-X '$(DTM_ROOT)/cmd/devstream/list.PluginsName=$(PLUGINS)'

FIND := find . -path './cmd/**/*.go' -o -path './test/**/*.go' -o -path './pkg/**/*.go' -o -path './internal/**/*.go'
GITHOOK := $(shell cp -f hack/githooks/* .git/hooks/)

# COLORS
RED    = $(shell printf "\33[31m")
GREEN  = $(shell printf "\33[32m")
WHITE  = $(shell printf "\33[37m")
YELLOW = $(shell printf "\33[33m")
RESET  = $(shell printf "\33[0m")


ifeq ($(GOOS),linux)
	MD5SUM=md5sum
else
	MD5SUM=md5 -q
endif

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(SELF_DIR) && pwd -P))
endif

ifeq ($(origin VERSION), undefined)
# the VERSION is a number, like 0.6.0
# it doesn't contain the prefix v, not v0.6.0, but 0.6.0
VERSION := $(shell git describe --tags --always --match='v*' | cut -c 2-)
endif

ifeq ($(origin PLUGINS_DIR),undefined)
PLUGINS_DIR := ${HOME}/.devstream/plugins
$(shell mkdir -p $(PLUGINS_DIR))
endif

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9.%-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Remove dtm and plugins. It's best to run a "clean" before "build".
	-rm -rf $(PLUGINS_DIR)
	-rm -f dtm*

.PHONY: build-core
build-core: generate fmt lint vet mod-tidy ## Build dtm core only, without plugins, locally.
	go build -trimpath -gcflags="all=-N -l" -ldflags "$(GO_LDFLAGS)" -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/
	@-rm -f dtm
	@cp dtm-${GOOS}-${GOARCH} dtm
	@echo "${GREEN}✔'dtm' has been generated in the current directory($(PWD))!${RESET}"

.PHONY: build-plugin.%
build-plugin.%: generate fmt lint vet mod-tidy ## Build one dtm plugin, like "make build-plugin.argocd".
	$(eval plugin_name := $(strip $*))
	@[ -d  $(ROOT_DIR)/cmd/plugin/$(plugin_name) ] || { echo -e "\n${RED}✘ Plugin '$(plugin_name)' not found!${RESET} The valid plugin name is as follows (Eg. You can use  ${YELLOW}make build-plugin.argocdapp${RESET} to build argocdapp plugin): \n\n$(shell ls ./cmd/plugin/)\n"; exit 1; }
	@echo "$(YELLOW)Building plugin '$(plugin_name)'$(RESET)"
	${GO_PLUGIN_BUILD} -o ${PLUGINS_DIR}/${plugin_name}-${PLUGIN_SUFFIX}.so ${ROOT_DIR}/cmd/plugin/${plugin_name}
	@$(MAKE) md5-plugin.$(plugin_name)

.PHONY: build-plugins
build-plugins: generate fmt vet mod-tidy $(addprefix build-plugin.,$(PLUGINS)) ## Build dtm plugins only. Use multi-threaded like "make build-plugins -j8" to speed up.

.PHONY: build
build: build-core build-plugins ## Build everything. Use multi-threaded like "make build -j8" to speed up.

.PHONY: md5
md5: md5-plugins ## Create md5 sums for all plugins.

.PHONY: md5-plugins
md5-plugins: $(addprefix md5-plugin.,$(PLUGINS))

.PHONY: md5-plugin.%
md5-plugin.%:
	$(eval plugin_name := $(strip $*))
	${MD5SUM} $(PLUGINS_DIR)/${plugin_name}-${PLUGIN_SUFFIX}.so > $(PLUGINS_DIR)/${plugin_name}-${PLUGIN_SUFFIX}.md5

.PHONY: fmt
fmt: verify.goimports ## Run 'go fmt' & goimports against code.
	@echo "$(YELLOW)Formating codes$(RESET)"
	@$(FIND) -type f | xargs gofmt -s -w
	@$(FIND) -type f | xargs ${GOPATH}/bin/goimports -w -local $(DTM_ROOT)
	@go mod edit -fmt

.PHONY: generate
generate: ## Run "go generate ./...".
	go generate ./...

.PHONY: lint
lint: verify.golangcilint ## Run 'golangci-lint' against code.
	@echo "$(YELLOW)Run golangci to lint source codes$(RESET)"
	@${GOPATH}/bin/golangci-lint -c $(ROOT_DIR)/.golangci.yml run $(ROOT_DIR)/...

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

.PHONY: verify.%
verify.%:
	@if ! command -v $* >/dev/null 2>&1; then $(MAKE) install.$*; fi

.PHONY: install.goimports
install.goimports:
	@go install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.golangcilint
install.golangcilint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
