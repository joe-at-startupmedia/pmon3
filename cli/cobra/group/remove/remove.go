package remove

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
)

var Cmd = &cobra.Command{
	Use:   "remove [group_name_or_id(s)] [process_name_or_id(s)]",
	Short: "Remove process(es) from group(s)",
	Args:  cobra.ExactArgs(2),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Remove(args[0], args[1])
	},
}
