package list

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller/base"
	"pmon3/cli/controller/group"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all groups",
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Show()
	},
}
