package cmd

import (
	"os"

	"github.com/devstream-io/devstream/internal/response"

	"github.com/devstream-io/devstream/internal/log"
	"github.com/devstream-io/devstream/internal/pkg/patch"

	"github.com/spf13/cobra"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "apply a diff file to an original",
	Long: `patch will take a patch file containing any of the four forms of difference listing
produced by the diff program and apply those differences to an original file,
producing a patched version. If patchfile is omitted, or is a hyphen,
the patch will be read from the standard input.

e.g.
- dtm patch file.patch
- dtm patch file.patch -ojson
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			errMsg := "Incorrect number of arguments"
			log.Error(errMsg)
			r := response.New(response.StatusError, response.MessageError, errMsg)
			r.Print(OutputFormat)
			os.Exit(1)
		}
		err := patch.Patch(args[0])
		if err != nil {
			log.Errorf("patch error: %v", err)
			r := response.New(response.StatusError, response.MessageError, err.Error())
			r.Print(OutputFormat)
		} else {
			r := response.New(response.StatusOK, response.MessageOK, "")
			r.Print(OutputFormat)
		}
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
}
