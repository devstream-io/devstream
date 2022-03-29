VERSION=0.3.0

MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH=$(patsubst %/,%,$(dir $(MKFILE_PATH)))/build/working_dir

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GO_BUILD=go build -buildmode=plugin -trimpath -gcflags="all=-N -l"

PLUGINS_CMD_ROOT=./cmd/plugin
PLUGINS_DIR=$(shell find ${PLUGINS_CMD_ROOT} -name "main.go" -exec dirname {} \;)
PLUGINS_NAME=$(notdir ${PLUGINS_DIR})
PLUGIN_SUFFIX=${GOOS}-${GOARCH}_${VERSION}

ifeq ($(GOOS),linux)
	MD5SUM=md5sum
else
	MD5SUM=md5 -q
endif

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9.%-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Remove dtm and plugins. It's best to run a "clean" before "build".
	rm -rf .devstream
	rm -f dtm*
	rm -f *.md5
	rm -rf build/working_dir

.PHONY: build-core
build-core: fmt vet mod-tidy ## Build dtm core only, without plugins, locally.
	go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/internal/pkg/version.Version=${VERSION}" -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/
	$(MAKE) md5-core
	rm -f dtm
	cp dtm-${GOOS}-${GOARCH} dtm

.PHONY: build-plugin.%
build-plugin.%: fmt vet mod-tidy mkdir.devstream ## Build one dtm plugin, like "make build-plugin.argocd"
	$(eval plugin_name := $(strip $*))
	${GO_BUILD} -o .devstream/${plugin_name}-${PLUGIN_SUFFIX}.so ${PLUGINS_CMD_ROOT}/${plugin_name}
	$(MAKE) md5-plugin.$(plugin_name)

.PHONY: build-plugins
build-plugins: fmt vet mod-tidy $(addprefix build-plugin.,$(PLUGINS_NAME)) ## Build dtm plugins only. Use multi-threaded like "make build-plugins -j8" to speed up.

.PHONY: build
build: build-core build-plugins ## Build everything. Use multi-threaded like "make build -j8" to speed up.

.PHONY: md5
md5: md5-core md5-plugins ## Create md5 sums for all plugins and dtm

.PHONY: md5-core
md5-core:
	${MD5SUM} dtm-${GOOS}-${GOARCH} > dtm-${GOOS}-${GOARCH}.md5

.PHONY: md5-plugins
md5-plugins: $(addprefix md5-plugin.,$(PLUGINS_NAME))

.PHONY: md5-plugin.%
md5-plugin.%:
	$(eval plugin_name := $(strip $*))
	${MD5SUM} .devstream/${plugin_name}-${PLUGIN_SUFFIX}.so > .devstream/${plugin_name}-${PLUGIN_SUFFIX}.md5


.PHONY: build-linux-amd64
build-linux-amd64: ## Cross-platform build for "linux/amd64".
	echo "Building in ${BUILD_PATH}"
	mkdir -p .devstream
	rm -rf ${BUILD_PATH} && mkdir ${BUILD_PATH}
	docker buildx build --platform linux/amd64 --load -t mericodev/stream-builder:v${VERSION} --build-arg VERSION=${VERSION} -f build/package/Dockerfile .
	cp -r go.mod go.sum cmd internal pkg build/package/build_linux_amd64.sh ${BUILD_PATH}/
	chmod +x ${BUILD_PATH}/build_linux_amd64.sh
	docker run --rm --platform linux/amd64 -v ${BUILD_PATH}:/devstream mericodev/stream-builder:v${VERSION}
	mv ${BUILD_PATH}/output/*.so .devstream/
	mv ${BUILD_PATH}/output/*.md5 .devstream/
	mv ${BUILD_PATH}/output/dtm* .
	rm -rf ${BUILD_PATH}

.PHONY: fmt
fmt: ## Run 'go fmt' & goimports against code.
	goimports -local="github.com/merico-dev/stream" -d -w cmd
	goimports -local="github.com/merico-dev/stream" -d -w pkg
	goimports -local="github.com/merico-dev/stream" -d -w internal
	goimports -local="github.com/merico-dev/stream" -d -w test
	go fmt ./...

.PHONY: vet
vet: ## Run "go vet ./...".
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Run "go mod tidy".
	go mod tidy

.PHONY: mkdir.devstream
mkdir.devstream:  ## Create ".devstream" (default directory for plugins) directory.
	mkdir -p .devstream

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
