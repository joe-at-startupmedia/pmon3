package initialize

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"time"

	"github.com/spf13/cobra"
)

var (
	appsConfigOnly bool
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "initialize all stopped processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()
	},
}

func init() {
	Cmd.Flags().BoolVarP(&appsConfigOnly, "apps-config-only", "c", false, "only initialize processes specified in the Apps Config file")
}

func Initialize() {
	base.OpenSender()
	defer base.CloseSender()
	arg1 := ""
	if appsConfigOnly {
		arg1 = "apps-config-only"
	}
	sent := base.SendCmd("init", arg1)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	list.Show()
}
