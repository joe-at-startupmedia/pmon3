package reset

import (
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/list"
	"pmon3/pmond/protos"
)

func Reset(idOrName string) *protos.CmdResp {

	sent := base.SendCmd("reset", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		list.Show()
	}
	return newCmdResp
}
