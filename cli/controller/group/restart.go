package group

import (
	"pmon3/cli"
	"pmon3/cli/controller"
	"pmon3/cli/controller/base"
	"pmon3/pmond/protos"
	"time"
)

func Restart(idOrName string, flags string) *protos.CmdResp {
	sent := base.SendCmdArg2("group_restart", idOrName, flags)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		controller.List()
	}
	return newCmdResp
}
