package drop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "drop",
	Short: "Delete all processes",
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		controller.Drop(forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill before deleting processes")
}
