package controller

import (
	model2 "pmon3/model"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Export(cmd *protos2.Cmd) *protos2.CmdResp {

	orderBy := cmd.GetArg1()

	var all []model2.Process
	var err error

	if len(orderBy) > 0 {
		all, err = repo.Process().FindAllOrdered(orderBy)
	} else {
		all, err = repo.Process().FindAll()
	}

	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	processConfig := model2.ProcessConfig{
		Processes: make([]model2.ExecFlags, len(all)),
	}
	for i, p := range all {
		processConfig.Processes[i] = *p.ToExecFlags()
	}

	newCmdResp := protos2.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: processConfig.Json(),
	}
	return &newCmdResp
}
