package controller

import (
	"fmt"
	"pmon3/model"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	"pmon3/protos"
)

func Kill(cmd *protos.Cmd) *protos.CmdResp {
	forced := cmd.GetArg1() == "force"
	return KillByParams(cmd, forced)
}

//KillByParams
/*
 * status param is the desired state to persist
 * this can either be status stopped or closed
 */
func KillByParams(cmd *protos.Cmd, forced bool) *protos.CmdResp {

	all, err := repo.Process().FindByStatus(model.StatusRunning)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("error finding running processes: %w", err))
	} else if len(all) == 0 {
		return base.ErroredCmdResp(cmd, fmt.Errorf("could not find running processes"))
	}

	for _, process := range all {
		_ = StopByParams(cmd, process.GetIdStr(), forced)
	}

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
