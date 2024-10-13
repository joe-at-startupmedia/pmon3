package list

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd"
	"pmon3/cli/cmd/base"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all processes",
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		cmd.List()
	},
}
