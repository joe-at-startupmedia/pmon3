package cmd

import (
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/protos"
	"time"
)

func Initialize(processConfigOnly bool, blocking bool) *protos.CmdResp {
	var (
		arg1 string
		arg2 string
	)

	if processConfigOnly {
		arg1 = "process-config-only"
	}
	if blocking {
		arg2 = "blocking"
	}

	sent := base.SendCmdArg2("init", arg1, arg2)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		time.Sleep(cli.Config.GetCmdExecResponseWait())
		List()
	}
	return newCmdResp
}
