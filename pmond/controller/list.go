package controller

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func List(cmd *protos.Cmd) *protos.CmdResp {
	var all []model.Process
	err := pmond.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 run err: %v", err)
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
