package stop

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/one"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "stop [id or name]",
	Short: "Stop a process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force the process to stop")
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("stop", args[0], "force")
	} else {
		sent = base.SendCmd("stop", args[0])
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(cli.Config.GetCmdExecResponseWait())
	p := model.ProcessFromProtobuf(newCmdResp.GetProcess())
	table_one.Render(p.RenderTable())
}
