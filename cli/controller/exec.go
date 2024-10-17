package controller

import (
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/model"
	"pmon3/protos"
	"time"
)

func Exec(file string, ef model.ExecFlags) *protos.CmdResp {
	ef.File = file
	sent := base.SendCmd("exec", ef.Json())
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		List()
	}
	return newCmdResp
}
