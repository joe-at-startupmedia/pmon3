package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func List(cmd *protos.Cmd) *protos.CmdResp {
	all, err := repo.Process().FindAllOrdered("ID")
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}
	newProcessList := protos.ProcessList{}
	for _, p := range all {
		newProcessList.Processes = append(newProcessList.Processes, p.ToProtobuf())
	}
	newCmdResp := protos.CmdResp{
		Id:          cmd.GetId(),
		Name:        cmd.GetName(),
		ProcessList: &newProcessList,
	}
	return &newCmdResp
}
