package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func List(cmd *protos2.Cmd) *protos2.CmdResp {
	all, err := repo.Process().FindAllOrdered("ID")
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}
	newProcessList := protos2.ProcessList{}
	for _, p := range all {
		newProcessList.Processes = append(newProcessList.Processes, p.ToProtobuf())
	}
	newCmdResp := protos2.CmdResp{
		Id:          cmd.GetId(),
		Name:        cmd.GetName(),
		ProcessList: &newProcessList,
	}
	return &newCmdResp
}
