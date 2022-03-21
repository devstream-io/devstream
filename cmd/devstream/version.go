package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/cmd/devstream/version"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of DevStream",
	Long:  `All software has versions. This is DevStream's`,
	Run:   versionCMDFunc,
}

func versionCMDFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Version:", version.Version, "\nSupported plugins MD5:", version.MD5String)
}
