package group

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/protos"
	"time"
)

func Delete(idOrName string) *protos.CmdResp {
	sent := base.SendCmd("group_del", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		Show()
	}
	return newCmdResp
}
