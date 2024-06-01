package initialize

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Restart all stopped processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()
	},
}

func Initialize() {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("init", "")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	list.Show()
}
