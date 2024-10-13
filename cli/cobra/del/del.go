package del

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/del"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "del [id or name]",
	Short: "Delete process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		del.Del(args[0], forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "kill the process before deletion")
}
