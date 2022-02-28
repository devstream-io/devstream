package version

// Version is the version of DevStream.
// Assignment by the command:
// `go build -trimpath -gcflags="all=-N -l" -ldflags "-X github.com/merico-dev/stream/cmd/devstream/version.Version=${VERSION}" \
// -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/`
// See the Makefile for more info.
var Version string
