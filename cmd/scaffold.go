package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/log"
	"github.com/devstream-io/devstream/internal/pkg/scaffold"
)

var structure string

// scaffoldCmd represents the scaffold command
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "scaffold is used to generate folder and file structure",
	Long: `
dtm scaffold "
project/
├── src/
│   ├── main.go
│   └── utils/
│       ├── file1.go
│       └── file2.go
└── README.md
"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Error("Incorrect number of arguments")
			os.Exit(1)
		}
		if err := scaffold.Scaffold(args[0]); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)

	scaffoldCmd.Flags().StringVarP(&structure, "structure", "s", "", "structure specify the folder and file structure")
}
