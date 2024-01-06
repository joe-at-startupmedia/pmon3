package stop

import (
	table_one "pmon3/cli/output/one"
	"pmon3/cli/pmq"
	"pmon3/pmond"
	"pmon3/pmond/model"
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
	pmq.New()
	if forceKill {
		pmq.SendCmdArg2("stop", args[0], "force")
	} else {
		pmq.SendCmd("stop", args[0])
	}
	newCmdResp := pmq.GetResponse()
	if len(newCmdResp.GetError()) > 0 {
		pmond.Log.Fatalf(newCmdResp.GetError())
	}
	time.Sleep(pmond.Config.GetCmdExecResponseWait())
	p := model.FromProtobuf(newCmdResp.GetProcess())
	table_one.Render(p.RenderTable())
	pmq.Close()
}
