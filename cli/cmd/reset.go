package cmd

import (
	"pmon3/cli/cmd/base"
	"pmon3/pmond/protos"
)

func Reset(idOrName string) *protos.CmdResp {

	sent := base.SendCmd("reset", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		List()
	}
	return newCmdResp
}
