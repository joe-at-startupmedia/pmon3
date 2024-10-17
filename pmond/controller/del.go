package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/del"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Delete(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return DeleteByParams(cmd, idOrName, forced)
}

func DeleteByParams(cmd *protos2.Cmd, idOrName string, forced bool) *protos2.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("process (%s) does not exist", idOrName))
	}

	err = del.ByProcess(p, forced)
	newCmdResp := protos2.CmdResp{
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
