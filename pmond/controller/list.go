package controller

import (
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func List(cmd *protos.Cmd) *protos.CmdResp {
	var all []model.Process
	err := db.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 can find processes: %v", err)
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
