package controller

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Export(cmd *protos.Cmd) *protos.CmdResp {
	all, err := repo.Process().FindAll()
	if err != nil {
		pmond.Log.Fatalf("pmon3 can find processes: %v", err)
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
