package reset

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	idOrNameFlag string
)

var Cmd = &cobra.Command{
	Use:   "reset [process_id_or_name]",
	Short: "Reset the restart counter(s)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		if len(args) == 1 {
			idOrNameFlag = args[0]
		}
		controller.Reset(idOrNameFlag)
	},
}
