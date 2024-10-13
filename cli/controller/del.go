package controller

import (
	"pmon3/cli/controller/base"
	"pmon3/cli/output/process/one"
	"pmon3/pmond/protos"

	"pmon3/pmond/model"
)

func Del(idOrName string, forceKill bool) *protos.CmdResp {
	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("del", idOrName, "force")
	} else {
		sent = base.SendCmd("del", idOrName)
	}
	newCmdResp := base.GetResponse(sent)
	process := newCmdResp.GetProcess()
	if process != nil {
		p := model.ProcessFromProtobuf(newCmdResp.GetProcess())
		table_one.Render(p.RenderTable())
	}

	return newCmdResp
}
