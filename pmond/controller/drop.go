package controller

import (
	"errors"
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Drop(cmd *protos.Cmd) *protos.CmdResp {
	forced := (cmd.GetArg1() == "force")
	return DropByParams(cmd, forced, model.StatusStopped)
}

func DropByParams(cmd *protos.Cmd, forced bool, status model.ProcessStatus) *protos.CmdResp {

	all, err := repo.Process().FindAll()
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("Error finding processes: %w", err))
	} else if len(all) == 0 {
		return base.ErroredCmdResp(cmd, errors.New("There are no processes"))
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
