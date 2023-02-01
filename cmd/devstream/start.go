package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/start"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var startCMD = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `start.`,
	Run:   startCMDFunc,
}

func startCMDFunc(_ *cobra.Command, _ []string) {
	err := start.Start()
	if err != nil {
		log.Fatal(err)
	}
}
