DTM_ROOT=github.com/devstream-io/devstream
SELF_DIR=$(dir $(lastword $(MAKEFILE_LIST)))
GOOS=$(shell go env GOOS)
GOPATH=$(shell go env GOPATH)
GOARCH=$(shell go env GOARCH)


# COLORS
RED    = $(shell printf "\33[31m")
GREEN  = $(shell printf "\33[32m")
WHITE  = $(shell printf "\33[37m")
YELLOW = $(shell printf "\33[33m")
RESET  = $(shell printf "\33[0m")

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(SELF_DIR) && pwd -P))
endif

FIND := find . -path './cmd/*.go' -o -path './internal/pkg/*.go'

.PHONY: build
build: fmt vet mod-tidy ## Build dtm core only, without plugins, locally.
	go build -trimpath -o dtm .
	@echo "${GREEN}✔'dtm' has been generated in the current directory($(PWD))!${RESET}"

.PHONY: fmt
fmt: verify.goimports ## Run 'go fmt' & goimports against code.
	@echo "$(YELLOW)Formating codes$(RESET)"
	@$(FIND) -type f | xargs gofmt -s -w
	@$(FIND) -type f | xargs ${GOPATH}/bin/goimports -w -local $(DTM_ROOT)
	@go mod edit -fmt

#.PHONY: lint
#lint: verify.golangcilint ## Run 'golangci-lint' against code.
#	@echo "$(YELLOW)Run golangci to lint source codes$(RESET)"
#	@${GOPATH}/bin/golangci-lint -c $(ROOT_DIR)/.golangci.yml run $(ROOT_DIR)/...

.PHONY: vet
vet: ## Run "go vet ./...".
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Run "go mod tidy".
	go mod tidy

.PHONY: verify.%
verify.%:
	@if ! command -v $* >/dev/null 2>&1; then $(MAKE) install.$*; fi

.PHONY: install.goimports
install.goimports:
	@go install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.golangcilint
install.golangcilint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
