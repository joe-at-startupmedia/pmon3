package export

import (
	"github.com/spf13/cobra"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
)

var (
	format  = "json"
	orderBy string
)

var Cmd = &cobra.Command{
	Use:   "export [format]",
	Short: "Export Process Configuration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		if len(args) == 1 {
			format = args[0]
		}
		controller.Export(format, orderBy)
	},
}

func init() {
	Cmd.Flags().StringVarP(&orderBy, "order", "o", "id", "the field by which to order by")
}
