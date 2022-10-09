package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PreRunFunc func(cmd *cobra.Command, args []string)

// BindPFlags accepts a list of flag names and binds the corresponding flags to Viper. This function servers
// as a workaround to a known bug in Viper, in which two commands can't share the same flag name.
// To learn more about the bug, see: https://github.com/spf13/viper/issues/233
func BindPFlags(flagNames []string) PreRunFunc {
	return func(cmd *cobra.Command, args []string) {
		for _, name := range flagNames {
			if err := viper.BindPFlag(name, cmd.Flags().Lookup(name)); err != nil {
				log.Fatalf("Failed to bind flag %s: %s", name, err)
			}
		}
	}
}
