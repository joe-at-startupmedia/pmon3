package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func ResetCounter(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()

	if len(idOrName) > 0 {
		p, err := repo.Process().FindByIdOrName(idOrName)
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("error finding processes: %s, %w", idOrName, err))
		}
		p.ResetRestartCount()
	} else {
		all, err := repo.Process().FindAll()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("error finding all processes: %v", err))
		}
		for _, p := range all {
			p.ResetRestartCount()
		}
	}

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
