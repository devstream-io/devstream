package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/upgrade"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var upgradeCMD = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade dtm to the latest release version",
	Long:  `Upgrade dtm to the latest release version.`,
	Run:   upgradeCMDFunc,
}

func upgradeCMDFunc(cmd *cobra.Command, args []string) {
	if err := upgrade.Upgrade(continueDirectly); err != nil {
		log.Fatal(err)
	}
}

func init() {
	addFlagContinueDirectly(upgradeCMD)
}
