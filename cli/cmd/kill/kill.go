package kill

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

var (
	forceKillFlag bool
)

var Cmd = &cobra.Command{
	Use:   "kill",
	Short: "Terminate all processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Kill(forceKillFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKillFlag, "force", "f", false, "force kill all processes")
}

func Kill(forceKill bool) *protos.CmdResp {
	var sent *protos.Cmd

	if forceKill {
		sent = base.SendCmd("kill", "force")
	} else {
		sent = base.SendCmd("kill", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
