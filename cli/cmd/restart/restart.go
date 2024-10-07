package restart

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
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
			Restart(cmd.CalledAs(), args[0], flag.Json())
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

func Restart(calledAs string, idOrName string, flags string) *protos.CmdResp {
	sent := base.SendCmdArg2(calledAs, idOrName, flags)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
