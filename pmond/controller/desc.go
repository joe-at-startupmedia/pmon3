package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	"pmon3/protos"
)

func Desc(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("process (%s) does not exist", idOrName))
	}
	newCmdResp := protos.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
