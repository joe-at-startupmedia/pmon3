package create

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/create"
)

var Cmd = &cobra.Command{
	Use:   "create [group_name]",
	Short: "Create a new group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		create.Create(args[0])
	},
}
