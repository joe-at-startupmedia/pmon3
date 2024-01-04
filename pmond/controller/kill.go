package controller

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Kill(cmd *protos.Cmd) *protos.CmdResp {
	forced := (cmd.GetArg1() == "force")
	return KillByParams(cmd, forced, model.StatusStopped)
}

/**
 * status param is the desired state to persist
 * this can either be status stopped or closed
 */
func KillByParams(cmd *protos.Cmd, forced bool, status model.ProcessStatus) *protos.CmdResp {

	var all []model.Process
	err := pmond.Db().Find(&all, "status = ?", model.StatusRunning).Error
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("Error finding running processes: %+v", err))
	} else if len(all) == 0 {
		return ErroredCmdResp(cmd, "Could not find running processes")
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
