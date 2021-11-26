GOPROXY ?= https://goproxy.io
GOOS ?= linux
GOARCH ?= amd64

mkfile_path=$(abspath $(lastword $(MAKEFILE_LIST)))
ROOT_PATH=$(dir $(mkfile_path))

build:
	go get ./...
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/githubactions_0.0.1.so ./cmd/githubactions/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocd_0.0.1.so ./cmd/argocd/
	go build -buildmode=plugin -trimpath -gcflags="all=-N -l" -o plugins/argocdapp_0.0.1.so ./cmd/argocdapp/
	go build -trimpath -gcflags="all=-N -l" -o dsm ./cmd/devstream/

# dsm and plugins will be located at build/output/
build-with-docker:
	rm -rf build/ && mkdir build
	# update the deps layer if needed
	docker build --build-arg GOPROXY=${GOPROXY} -t devstream_build:v0.1 .
	cp -r go.mod go.sum cmd internal build.sh build/
	docker run -v ${ROOT_PATH}build:/devstream -e GOOS=${GOOS} -e GOARCH=${GOARCH} --rm devstream_build:v0.1
