package restart

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group"
	"pmon3/pmond/model"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:   "restart [group_id_or_name]",
	Short: "(Re)start processes by group id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cobraCommand *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		group.Restart(args[0], flag.Json())
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the processes stdout log dir")
}
