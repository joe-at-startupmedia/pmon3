package kill

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/kill"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "kill",
	Short: "Terminate all processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		kill.Kill(forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill all processes")
}
