package group

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/protos"
	"time"
)

func Create(groupName string) *protos.CmdResp {
	sent := base.SendCmd("group_create", groupName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		Show()
	}
	return newCmdResp
}
