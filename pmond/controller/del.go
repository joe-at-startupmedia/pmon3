package controller

import (
	"fmt"
	"pmon3/pmond/controller/base/del"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Delete(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return DeleteByParams(cmd, idOrName, forced)
}

func DeleteByParams(cmd *protos.Cmd, idOrName string, forced bool) *protos.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		newCmdResp := protos.CmdResp{
			Id:    cmd.GetId(),
			Name:  cmd.GetName(),
			Error: fmt.Sprintf("Process (%s) does not exist", idOrName),
		}
		return &newCmdResp
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
