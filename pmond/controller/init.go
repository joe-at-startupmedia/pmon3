package controller

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Initialize(cmd *protos.Cmd) *protos.CmdResp {
	var all []model.Process
	err := pmond.Db().Find(&all, "status = ?", model.StatusStopped).Error
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("Error finding stopped processes: %w", err))
	} else if len(all) == 0 {
		return ErroredCmdResp(cmd, fmt.Errorf("Could not find any stopped processes"))
	}

	var (
		cr       *protos.CmdResp
		hasError bool = false
	)

	for _, process := range all {
		pmond.Log.Debugf("restart: %s", process.GetIdStr())
		cr := RestartByParams(cmd, process.GetIdStr(), "{}", false)
		if len(cr.GetError()) > 0 {
			pmond.Log.Debugf("encountered error attempting to restart: %s", cr.GetError())
			hasError = true
			break
		}
	}

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	if hasError {
		newCmdResp.Error = cr.GetError()
	}
	return &newCmdResp
}
