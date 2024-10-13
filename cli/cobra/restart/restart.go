package restart

import (
	"github.com/spf13/cobra"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/restart"
	"pmon3/pmond/model"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:     "restart [id or name]",
	Short:   "(Re)start a process by id or name",
	Aliases: []string{"start"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(flag.User) > 0 && flag.User == "root" && !base.IsRoot() {
			base.OutputError("cannot set process user to root without sudo")
		} else {
			base.OpenSender()
			defer base.CloseSender()
			restart.Restart(cmd.CalledAs(), args[0], flag.Json())
		}
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVar(&flag.LogDir, "log_dir", "", "the processes stdout log dir")
	Cmd.Flags().StringSliceVarP(&flag.Dependencies, "dependencies", "d", []string{}, "provide a list of process names this process depends on")
}
