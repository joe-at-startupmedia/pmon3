package controller

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
)

func Stop(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := (cmd.GetArg2() == "force")
	return StopByParams(cmd, idOrName, forced, model.StatusStopped)
}

func StopByParams(cmd *protos.Cmd, idOrName string, forced bool, status model.ProcessStatus) *protos.CmdResp {
	err, p := model.FindProcessByIdOrName(pmond.Db(), idOrName)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("could not find process: %+v", err))
	}

	// check process is running
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	//if process is not currently running
	if os.IsNotExist(err) {
		if p.Status != status {
			p.Status = status
			if err := pmond.Db().Save(&p).Error; err != nil {
				return ErroredCmdResp(cmd, fmt.Sprintf("stop process error: %+v", err))
			}
		}
	}

	// try to kill the process
	err = process.SendOsKillSignal(p, status, forced)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("stop process error: %+v", err))
	}

	pmond.Log.Infof("stop process %s success \n", p.Stringify())

	p.ResetRestartCount()

	newCmdResp := protos.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
