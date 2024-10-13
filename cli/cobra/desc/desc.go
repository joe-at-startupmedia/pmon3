package desc

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "List process information by id or name",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		controller.Desc(args[0])
	},
}
