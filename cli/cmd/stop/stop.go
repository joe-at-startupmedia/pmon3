package stop

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/one"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"
)

func Stop(idOrName string, forceKill bool) *protos.CmdResp {
	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("stop", idOrName, "force")
	} else {
		sent = base.SendCmd("stop", idOrName)
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		p := model.ProcessFromProtobuf(newCmdResp.GetProcess())
		table_one.Render(p.RenderTable())
	}
	return newCmdResp
}
