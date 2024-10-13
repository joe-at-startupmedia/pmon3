package kill

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "kill",
	Short: "Terminate all processes",
	Args:  cobra.NoArgs,
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		controller.Kill(forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill all processes")
}
