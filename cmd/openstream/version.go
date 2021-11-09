package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ironcore864/openstream/internal/pkg/version"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of openstream",
	Long:  `All software has versions. This is openstream's`,
	Run:   versionCMDFunc,
}

func versionCMDFunc(cmd *cobra.Command, args []string) {
	fmt.Println(version.VERSION)
}
