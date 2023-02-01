package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/create"
)

var createCMD = &cobra.Command{
	Use:   "create",
	Short: "create",
	Long:  `create.`,
	Run:   createCMDFunc,
}

func createCMDFunc(cmd *cobra.Command, args []string) {
	err := create.Create()
	if err != nil {
		log.Panic(err)
	}
}
