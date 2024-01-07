package kill

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
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
	if forceKill {
		base.SendCmd("kill", "force")
	} else {
		base.SendCmd("kill", "")
	}
	newCmdResp := base.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	//list command will call pmq.Close
	list.Show()
}
