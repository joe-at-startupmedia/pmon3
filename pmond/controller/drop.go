package controller

import (
	"errors"
	"fmt"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Drop(cmd *protos2.Cmd) *protos2.CmdResp {
	forced := cmd.GetArg1() == "force"
	return DropByParams(cmd, forced)
}

func DropByParams(cmd *protos2.Cmd, forced bool) *protos2.CmdResp {

	all, err := repo.Process().FindAll()
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("error finding processes: %w", err))
	} else if len(all) == 0 {
		return base.ErroredCmdResp(cmd, errors.New("there are no processes"))
	}

	for _, process := range all {
		_ = DeleteByParams(cmd, process.GetIdStr(), forced)
	}

	newCmdResp := protos2.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	return &newCmdResp
}
