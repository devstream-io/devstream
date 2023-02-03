package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/start"
)

var startCMD = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `start.`,
	Run:   startCMDFunc,
}

func startCMDFunc(_ *cobra.Command, _ []string) {
	err := start.Start()
	if err != nil && err.Error() != "^C" {
		fmt.Printf("Failed with error: %s", err)
	}
}
