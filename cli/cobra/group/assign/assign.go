package assign

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
)

var Cmd = &cobra.Command{
	Use:   "assign [group_name_or_ids(s)] [process_name_or_id(s)]",
	Short: "Assign group(s) to process(es)",
	Args:  cobra.ExactArgs(2),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Assign(args[0], args[1])
	},
}
