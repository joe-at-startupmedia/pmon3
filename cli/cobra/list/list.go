package list

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all processes",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		list.Show()
	},
}
