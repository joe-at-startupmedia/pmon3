package del

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
)

var Cmd = &cobra.Command{
	Use:   "del [group_name]",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Delete(args[0])
	},
}
