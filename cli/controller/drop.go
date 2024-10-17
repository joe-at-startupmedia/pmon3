package controller

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	protos2 "pmon3/protos"
	"time"
)

func Drop(forceKill bool) *protos2.CmdResp {
	var sent *protos2.Cmd

	if forceKill {
		sent = base.SendCmd("drop", "force")
	} else {
		sent = base.SendCmd("drop", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		List()
	}
	return newCmdResp
}
