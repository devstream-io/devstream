package version

import "regexp"

// Version is the version of DevStream.
// Assign the value when building with the -X parameter. Example:
// -X github.com/devstream-io/devstream/internal/pkg/version.Version=${VERSION}
// See the Makefile for more info.
var Version string

var Dev bool

func init() {
	// check is in dev version
	Dev = !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(Version)
}
