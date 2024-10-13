package restart

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
	"time"
)

func Restart(calledAs string, idOrName string, flags string) *protos.CmdResp {
	sent := base.SendCmdArg2(calledAs, idOrName, flags)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
