package options

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/spf13/cobra"
)

// to pre validate args of commands

// WithValidators returns a cobra RunFunc with the given validators
func WithValidators(runFunc func(cmd *cobra.Command, args []string), validators ...func(args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		for _, validator := range validators {
			if err := validator(args); err != nil {
				log.Fatal(err)
			}
		}
		runFunc(cmd, args)
	}

}

// ArgsCountEqual returns a validator which check the count of args
func ArgsCountEqual(count int) func(args []string) error {
	return func(args []string) error {
		if len(args) != count {
			return fmt.Errorf("illegal args count (expect %d, got %d)", count, len(args))
		}
		return nil
	}
}
