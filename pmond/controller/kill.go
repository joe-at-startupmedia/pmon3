package controller

import (
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Kill(cmd *protos.Cmd) *protos.CmdResp {
	forced := cmd.GetArg1() == "force"
	return KillByParams(cmd, forced, model.StatusStopped)
}

//KillByParams
/*
 * status param is the desired state to persist
 * this can either be status stopped or closed
 */
func KillByParams(cmd *protos.Cmd, forced bool, status model.ProcessStatus) *protos.CmdResp {

	all, err := repo.Process().FindByStatus(model.StatusRunning)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("error finding running processes: %w", err))
	} else if len(all) == 0 {
		return base.ErroredCmdResp(cmd, fmt.Errorf("could not find running processes"))
	}

	for _, process := range all {
		_ = StopByParams(cmd, process.GetIdStr(), forced, status)
	}

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
