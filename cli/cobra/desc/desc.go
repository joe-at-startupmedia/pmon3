package desc

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/desc"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "Show process information by id or name",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		desc.Desc(args[0])
	},
}
