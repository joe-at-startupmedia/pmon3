package export

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/export"
)

var (
	format  string
	orderBy string
)

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Process Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		export.Export(format, orderBy)
	},
}

func init() {
	Cmd.Flags().StringVarP(&format, "format", "f", "json", "the format to export")
	Cmd.Flags().StringVarP(&orderBy, "order", "o", "json", "the field by which to order by")
}
