package exec

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"
)

func Exec(file string, ef model.ExecFlags) *protos.CmdResp {
	ef.File = file
	sent := base.SendCmd("exec", ef.Json())
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		list.Show()
	}
	return newCmdResp
}
