package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
)

func completionCMD(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion",
		Short:                 "Generate the autocompletion script for dtm for the specified shell",
		Long:                  "See each sub-command's help for details on how to use the generated script.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactValidArgs(1),
	}

	binaryName := filepath.Base(os.Args[0])
	bash := &cobra.Command{
		Use:     "bash",
		Short:   "generate autocompletion script for bash",
		Example: completion.BashExample(binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return completion.CompletionBash(out, cmd)
		},
	}

	zsh := &cobra.Command{
		Use:     "zsh",
		Short:   "generate autocompletion script for zsh",
		Example: completion.ZshExample(binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return completion.CompletionZsh(out, cmd)
		},
	}

	fish := &cobra.Command{
		Use:     "fish",
		Short:   "generate autocompletion script for fish",
		Example: completion.FishExample(binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenFishCompletion(out, true)
		},
	}

	powershell := &cobra.Command{
		Use:     "powershell",
		Short:   "generate autocompletion script for powershell",
		Example: completion.PowershellExample(binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenPowerShellCompletionWithDesc(out)
		},
	}
	cmd.AddCommand(bash, zsh, fish, powershell)

	return cmd
}
