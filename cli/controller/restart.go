package controller

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/protos"
	"time"
)

func Restart(calledAs string, idOrName string, flags string) *protos.CmdResp {
	sent := base.SendCmdArg2(calledAs, idOrName, flags)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		List()
	}
	return newCmdResp
}
