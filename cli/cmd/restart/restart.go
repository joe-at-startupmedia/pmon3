package restart

import (
	"pmon3/cli/cmd/list"
	"pmon3/cli/pmq"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"time"

	"github.com/spf13/cobra"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:   "restart [id or name]",
	Short: "Restart a process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the processes stdout log dir")
}

func cmdRun(args []string, flags string) {
	pmq.New()
	pmq.SendCmdArg2("restart", args[0], flags)
	newCmdResp := pmq.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		pmond.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(pmond.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}
