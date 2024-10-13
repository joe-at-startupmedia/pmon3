package desc

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/desc"
)

var Cmd = &cobra.Command{
	Use:     "desc [id or name]",
	Aliases: []string{"show"},
	Short:   "Show group details and associated processes",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		desc.Desc(args[0])
	},
}
