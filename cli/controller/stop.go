package controller

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/cli/output/process/one"
	"pmon3/model"
	protos2 "pmon3/protos"
	"time"
)

func Stop(idOrName string, forceKill bool) *protos2.CmdResp {
	var sent *protos2.Cmd
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
