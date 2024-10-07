package initialize

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

var (
	processConfigOnlyFlag bool
	blockingFlag          bool
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize all stopped processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Initialize(processConfigOnlyFlag, blockingFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
	Cmd.Flags().BoolVarP(&blockingFlag, "blocking", "b", false, "return a response only after all processes have been queued")
}

func Initialize(processConfigOnly bool, blocking bool) *protos.CmdResp {
	var (
		arg1 string
		arg2 string
	)

	if processConfigOnly {
		arg1 = "process-config-only"
	}
	if blocking {
		arg2 = "blocking"
	}

	sent := base.SendCmdArg2("init", arg1, arg2)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
