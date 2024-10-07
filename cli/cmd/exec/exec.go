package exec

import (
	"os/user"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

// process failed auto restart
var flags model.ExecFlags

var Cmd = &cobra.Command{
	Use:     "exec [application_binary]",
	Aliases: []string{"run"},
	Short:   "Spawn a new process",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(flags.User) > 0 && flags.User == "root" && !base.IsRoot() {
			base.OutputError("cannot set process user to root without sudo")
			return
		} else if flags.User == "" {
			user, err := user.Current()
			if err == nil {
				flags.User = user.Username
			}
		}
		base.OpenSender()
		defer base.CloseSender()
		Exec(args[0], flags)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flags.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flags.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flags.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flags.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVarP(&flags.EnvVars, "env-vars", "e", "", "the processes environment variables (space-delimited)")
	Cmd.Flags().StringVar(&flags.Name, "name", "", "the processes name")
	Cmd.Flags().StringVar(&flags.LogDir, "log-dir", "", "the processes stdout log dir")
	Cmd.Flags().StringSliceVarP(&flags.Dependencies, "dependencies", "d", []string{}, "provide a list of process names this process depends on")
	Cmd.Flags().StringSliceVarP(&flags.Groups, "groups", "g", []string{}, "provide a list of group names this process is associated to")
}

func Exec(file string, ef model.ExecFlags) *protos.CmdResp {
	ef.File = file
	sent := base.SendCmd("exec", ef.Json())
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
