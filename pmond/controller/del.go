package controller

import (
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Delete(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := (cmd.GetArg2() == "force")
	return DeleteByParams(cmd, idOrName, forced)
}

func DeleteByParams(cmd *protos.Cmd, idOrName string, forced bool) *protos.CmdResp {
	stopCmdResp := StopByParams(cmd, idOrName, forced, model.StatusStopped)
	if len(stopCmdResp.GetError()) > 0 {
		return ErroredCmdResp(cmd, stopCmdResp.GetError())
	}
	process := model.FromProtobuf(stopCmdResp.GetProcess())
	err := pmond.Db().Delete(process).Error
	_ = os.Remove(process.Log)
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	if process != nil {
		newCmdResp.Process = stopCmdResp.GetProcess()
	}
	if err != nil {
		newCmdResp.Error = err.Error()
	}
	return &newCmdResp
}
