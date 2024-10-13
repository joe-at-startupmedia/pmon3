package group

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/pmond/protos"
	"time"
)

func Remove(groupNameOrId string, processNameOrId string) *protos.CmdResp {
	sent := base.SendCmdArg2("group_remove", groupNameOrId, processNameOrId)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		Desc(groupNameOrId)
	}
	return newCmdResp
}
