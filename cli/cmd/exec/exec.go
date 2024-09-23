package exec

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
	"time"

	"github.com/spf13/cobra"
)

// process failed auto restart
var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:     "exec [application_binary]",
	Aliases: []string{"run"},
	Short:   "Spawn a new process",
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flag.SetCurrentUser()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVarP(&flag.EnvVars, "env-vars", "e", "", "the processes environment variables (space-delimited)")
	Cmd.Flags().StringVar(&flag.Name, "name", "", "the processes name")
	Cmd.Flags().StringVar(&flag.LogDir, "log-dir", "", "the processes stdout log dir")
	Cmd.Flags().StringSliceVarP(&flag.Dependencies, "dependencies", "d", []string{}, "provide a list of process names this process depends on")
	Cmd.Flags().StringSliceVarP(&flag.Groups, "groups", "g", []string{}, "provide a list of group names this process is associated to")
}

func cmdRun(args []string, flags string) {
	base.OpenSender()
	sent := base.SendCmdArg2("exec", args[0], flags)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}
