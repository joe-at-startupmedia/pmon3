package stop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "stop [id or name]",
	Short: "Stop a process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		cmd.Stop(args[0], forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force the process to stop")
}
