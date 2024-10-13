package del

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/del"
)

var Cmd = &cobra.Command{
	Use:   "del [group_name]",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		del.Delete(args[0])
	},
}
