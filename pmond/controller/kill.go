package controller

import (
	"fmt"
	"pmon3/model"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Kill(cmd *protos2.Cmd) *protos2.CmdResp {
	forced := cmd.GetArg1() == "force"
	return KillByParams(cmd, forced)
}

//KillByParams
/*
 * status param is the desired state to persist
 * this can either be status stopped or closed
 */
func KillByParams(cmd *protos2.Cmd, forced bool) *protos2.CmdResp {

	all, err := repo.Process().FindByStatus(model.StatusRunning)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("error finding running processes: %w", err))
	} else if len(all) == 0 {
		return base.ErroredCmdResp(cmd, fmt.Errorf("could not find running processes"))
	}

	for _, process := range all {
		_ = StopByParams(cmd, process.GetIdStr(), forced)
	}

	newCmdResp := protos2.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
