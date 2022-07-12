package main

import (
	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/upgrade"
)

var upgradeCMD = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade dtm to the latest version",
	Long:  `Upgrade dtm to the latest version.`,
	Run:   upgradeCMDFunc,
}

func upgradeCMDFunc(cmd *cobra.Command, args []string) {
	if err := upgrade.Upgrade(continueDirectly); err != nil {
		log.Fatal(err)
	}
}

func init() {
	upgradeCMD.Flags().BoolVarP(&continueDirectly, "yes", "y", false, "apply directly without confirmation")
}
