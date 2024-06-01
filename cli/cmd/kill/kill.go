package kill

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "kill",
	Short: "Terminate all processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Kill(model.StatusStopped)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force kill all processes")
}

func Kill(processStatus model.ProcessStatus) {
	base.OpenSender()
	defer base.CloseSender()

	var sent *protos.Cmd

	if forceKill {
		sent = base.SendCmd("kill", "force")
	} else {
		sent = base.SendCmd("kill", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	list.Show()
}
