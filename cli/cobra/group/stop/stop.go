package stop

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
)

var Cmd = &cobra.Command{
	Use:     "stop [group_id_or_name]",
	Aliases: []string{"show"},
	Short:   "Stop all processes associated to a group",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Stop(args[0])
	},
}
