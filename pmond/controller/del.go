package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/del"
	"pmon3/pmond/repo"
	"pmon3/protos"
)

func Delete(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return DeleteByParams(cmd, idOrName, forced)
}

func DeleteByParams(cmd *protos.Cmd, idOrName string, forced bool) *protos.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("process (%s) does not exist", idOrName))
	}

	err = del.ByProcess(p, forced)
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	if p != nil {
		newCmdResp.Process = p.ToProtobuf()
	}
	if err != nil {
		newCmdResp.Error = err.Error()
	}
	return &newCmdResp
}
