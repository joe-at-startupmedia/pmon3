package controller

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"time"
)

func Stop(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return StopByParams(cmd, idOrName, forced, model.StatusStopped)
}

func StopByParams(cmd *protos.Cmd, idOrName string, forced bool, status model.ProcessStatus) *protos.CmdResp {
	err, p := model.FindProcessByIdOrName(db.Db(), idOrName)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("could not find process: %w", err))
	}

	// check process is running
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	//if process is not currently running
	if os.IsNotExist(err) {
		if p.Status != status {
			p.Status = status
			if err := db.Db().Save(&p).Error; err != nil {
				return ErroredCmdResp(cmd, fmt.Errorf("stop process error: %w", err))
			}
		}
	}

	//we need to wait for the process to save before killing it to avoid a restart race condition
	time.Sleep(200 * time.Millisecond)
	// try to kill the process
	err = process.SendOsKillSignal(p, status, forced)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("stop process error: %w", err))
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
