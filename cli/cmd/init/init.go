package initialize

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"time"

	"github.com/spf13/cobra"
)

var (
	processConfigOnly bool
	blocking          bool
	arg1              string
	arg2              string
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize all stopped processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnly, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
	Cmd.Flags().BoolVarP(&blocking, "blocking", "b", false, "return a response only after all processes have been queued")
}

func Initialize() {
	base.OpenSender()
	defer base.CloseSender()
	if processConfigOnly {
		arg1 = "process-config-only"
	}
	if blocking {
		arg2 = "blocking"
	}
	sent := base.SendCmdArg2("init", arg1, arg2)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	list.Show()
}
