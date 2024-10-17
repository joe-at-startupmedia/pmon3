package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func ResetCounter(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()

	if len(idOrName) > 0 {
		p, err := repo.Process().FindByIdOrName(idOrName)
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("process (%s) does not exist", idOrName))
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

	newCmdResp := protos2.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
