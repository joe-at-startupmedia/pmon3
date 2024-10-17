package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Stop(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return StopByParams(cmd, idOrName, forced)
}

func StopByParams(cmd *protos2.Cmd, idOrName string, forced bool) *protos2.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	err = stop.ByProcess(p, forced)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	newCmdResp := protos2.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
