package controller

import (
	"errors"
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Drop(cmd *protos.Cmd) *protos.CmdResp {
	forced := (cmd.GetArg1() == "force")
	return DropByParams(cmd, forced, model.StatusStopped)
}

func DropByParams(cmd *protos.Cmd, forced bool, status model.ProcessStatus) *protos.CmdResp {

	var all []model.Process
	err := pmond.Db().Find(&all).Error
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("Error finding processes: %w", err))
	} else if len(all) == 0 {
		return ErroredCmdResp(cmd, errors.New("There are no processes"))
	}

	for _, process := range all {
		_ = DeleteByParams(cmd, process.GetIdStr(), forced)
	}

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
