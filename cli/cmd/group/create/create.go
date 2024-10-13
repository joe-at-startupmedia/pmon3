package create

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/group/list"
	"pmon3/pmond/protos"
	"time"
)

func Create(groupName string) *protos.CmdResp {
	sent := base.SendCmd("group_create", groupName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
