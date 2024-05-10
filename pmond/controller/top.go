package controller

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Top(cmd *protos.Cmd) *protos.CmdResp {
	var all []model.Process
	err := db.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 can find processes: %v", err)
	}
	pidsCsv := fmt.Sprintf("%d", os.Getpid())
	for _, p := range all {
		pidsCsv = fmt.Sprintf("%d,%s", p.Pid, pidsCsv)
	}
	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: pidsCsv,
	}
	return &newCmdResp
}
