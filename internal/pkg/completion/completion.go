package completion

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func FlagPluginsCompletion(cmd *cobra.Command, flag string) {
	if err := cmd.RegisterFlagCompletionFunc(flag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return list.PluginsNameSlice(), cobra.ShellCompDirectiveDefault
	}); err != nil {
		log.Warn(err)
	}
}

func FlagFilenameCompletion(cmd *cobra.Command, flagName string) {

	// Ref: https://github.com/spf13/cobra/blob/master/shell_completions.md#specify-valid-filename-extensions-for-flags-that-take-a-filename
	if err := cmd.MarkFlagFilename(flagName, "yaml", "yml"); err != nil {
		log.Warn(err)
	}
}

func FlagDirnameCompletion(cmd *cobra.Command, flagName string) {

	// Ref: https://github.com/spf13/cobra/blob/master/shell_completions.md#limit-flag-completions-to-directory-names
	if err := cmd.MarkFlagDirname(flagName); err != nil {
		log.Warn(err)
	}
}

func CompletionBash(out io.Writer, cmd *cobra.Command) error {
	err := cmd.Root().GenBashCompletion(out)

	// The default binary name downloaded from the Releases page is dtm-{os}-amd64
	// solve the problem that autocompletion fails when the name of the binary is not dtm
	if binary := filepath.Base(os.Args[0]); binary != "dtm" {
		renamedBinary := `
# the user renamed the dtm binary
if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_dtm %[1]s
else
    complete -o default -o nospace -F __start_dtm %[1]s
fi
`
		fmt.Fprintf(out, renamedBinary, binary)
	}

	return err
}

func CompletionZsh(out io.Writer, cmd *cobra.Command) error {
	err := cmd.Root().GenZshCompletionNoDesc(out)

	// The default binary name downloaded from the Releases page is dtm-{os}-amd64
	// solve the problem that autocompletion fails when the name of the binary is not dtm
	if binary := filepath.Base(os.Args[0]); binary != "dtm" {
		renamedBinary := `
# the user renamed the dtm binary
compdef _dtm %[1]s
`
		fmt.Fprintf(out, renamedBinary, binary)
	}

	fmt.Fprintf(out, "compdef _dtm dtm")

	return err
}

func BashExample(binary string) string {
	return fmt.Sprintf(`Load is completions in the current shell session:
# source <(%s completion bash)

Load completions for every new session:
(in Linux)# %s completion bash > /etc/bash_completion.d/dtm
(in MacOS)# %s completion bash > $(brew --prefix)/etc/bash_completion.d/dtm`, binary, binary, binary)
}

func ZshExample(binary string) string {
	return fmt.Sprintf(`Load is completions in the current shell session:
# source <(%s completion zsh)

Load completions for every new session:
# %s completion zsh > "${fpath[1]}/_dtm"`, binary, binary)
}

func FishExample(binary string) string {
	return fmt.Sprintf(`Load is completions in the current shell session:
#  %s completion fish | source

Load completions for every new session:
# %s completion fish > ~/.config/fish/completions/dtm.fish`, binary, binary)
}

func PowershellExample(binary string) string {
	return fmt.Sprintf(`Load is completions in the current shell session:
C:\> %s completion powershell | Out-String | Invoke-Expression

Load completions for every new session:
add the output of the above command to powershell profile.`, binary)
}
