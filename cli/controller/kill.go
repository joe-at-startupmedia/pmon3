package controller

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	protos2 "pmon3/protos"
	"time"
)

func Kill(forceKill bool) *protos2.CmdResp {
	var sent *protos2.Cmd

	if forceKill {
		sent = base.SendCmd("kill", "force")
	} else {
		sent = base.SendCmd("kill", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		List()
	}
	return newCmdResp
}
