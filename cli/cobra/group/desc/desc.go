package desc

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller/base"
	"pmon3/cli/controller/group"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "Show group details and associated processes",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Desc(args[0])
	},
}
