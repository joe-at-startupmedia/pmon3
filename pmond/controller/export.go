package controller

import (
	"pmon3/model"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	"pmon3/protos"
)

func Export(cmd *protos.Cmd) *protos.CmdResp {

	orderBy := cmd.GetArg1()

	var all []model.Process
	var err error

	if len(orderBy) > 0 {
		all, err = repo.Process().FindAllOrdered(orderBy)
	} else {
		all, err = repo.Process().FindAll()
	}

	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	processConfig := model.ProcessConfig{
		Processes: make([]model.ExecFlags, len(all)),
	}
	for i, p := range all {
		processConfig.Processes[i] = *p.ToExecFlags()
	}

	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: processConfig.Json(),
	}
	return &newCmdResp
}
