package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/version"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of devstream",
	Long:  `All software has versions. This is devstream's`,
	Run:   versionCMDFunc,
}

func versionCMDFunc(cmd *cobra.Command, args []string) {
	fmt.Println(version.VERSION)
}
