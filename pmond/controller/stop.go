package controller

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"time"
)

func Stop(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return StopByParams(cmd, idOrName, forced, model.StatusStopped)
}

func StopByParams(cmd *protos.Cmd, idOrName string, forced bool, status model.ProcessStatus) *protos.CmdResp {
	p, err := repo.Process().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	// check process is running
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	//if process is not currently running
	if os.IsNotExist(err) {
		if p.Status != status {
			if err := repo.ProcessOf(p).UpdateStatus(status); err != nil {
				return base.ErroredCmdResp(cmd, fmt.Errorf("stop process error: %w", err))
			}
		}
	}

	//we need to wait for the process to save before killing it to avoid a restart race condition
	time.Sleep(200 * time.Millisecond)
	// try to kill the process
	err = process.SendOsKillSignal(p, status, forced)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("stop process error: %w", err))
	}

	pmond.Log.Infof("stop process %s success", p.Stringify())

	p.ResetRestartCount()

	newCmdResp := protos.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
