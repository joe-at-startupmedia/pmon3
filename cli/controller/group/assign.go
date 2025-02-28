package group

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/protos"
	"time"
)

func Assign(groupNameOrId string, processNameOrId string) *protos.CmdResp {
	sent := base.SendCmdArg2("group_assign", groupNameOrId, processNameOrId)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		Desc(groupNameOrId)
	}
	return newCmdResp
}
