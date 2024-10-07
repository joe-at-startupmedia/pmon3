package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Stop(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return StopByParams(cmd, idOrName, forced)
}

func StopByParams(cmd *protos.Cmd, idOrName string, forced bool) *protos.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	err = stop.ByProcess(p, forced)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	newCmdResp := protos.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
