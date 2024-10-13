package drop

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
	"time"
)

func Drop(forceKill bool) *protos.CmdResp {
	var sent *protos.Cmd

	if forceKill {
		sent = base.SendCmd("drop", "force")
	} else {
		sent = base.SendCmd("drop", "")
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
