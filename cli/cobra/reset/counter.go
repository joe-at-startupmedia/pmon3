package reset

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
)

var (
	idOrNameFlag string
)

var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the restart counter(s)",
	Args:  cobra.ExactArgs(0),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		cmd.Reset(idOrNameFlag)
	},
}

func init() {
	Cmd.Flags().StringVarP(&idOrNameFlag, "process", "p", "", "the id or name of the process")
}
