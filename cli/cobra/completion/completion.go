package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                   "completion [bash|zsh]",
	Short:                 "Generate completion script",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		}
	},
}
