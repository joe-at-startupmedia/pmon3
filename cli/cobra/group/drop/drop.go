package drop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:     "drop [group_id_or_name]",
	Aliases: []string{"show"},
	Short:   "Delete all processes associated to a group",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Drop(args[0], forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill before deleting processes")
}
