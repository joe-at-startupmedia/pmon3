package reset

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/reset"
)

var (
	idOrNameFlag string
)

var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the restart counter(s)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		reset.Reset(idOrNameFlag)
	},
}

func init() {
	Cmd.Flags().StringVarP(&idOrNameFlag, "process", "p", "", "the id or name of the process")
}
