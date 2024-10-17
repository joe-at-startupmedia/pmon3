package controller

import (
	"fmt"
	"os"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Top(cmd *protos2.Cmd) *protos2.CmdResp {
	all, err := repo.Process().FindAll()
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}
	pidsCsv := fmt.Sprintf("%d", os.Getpid())
	for _, p := range all {
		pidsCsv = fmt.Sprintf("%d,%s", p.Pid, pidsCsv)
	}
	newCmdResp := protos2.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: pidsCsv,
	}
	return &newCmdResp
}
