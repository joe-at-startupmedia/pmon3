package controller

import (
	"fmt"
	"os"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Top(cmd *protos.Cmd) *protos.CmdResp {
	all, err := repo.Process().FindAll()
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
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
