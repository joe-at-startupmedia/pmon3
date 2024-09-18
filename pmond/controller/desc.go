package controller

import (
	"fmt"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Desc(cmd *protos.Cmd) *protos.CmdResp {
	val := cmd.GetArg1()
	p, err := repo.Process().FindByIdOrNameWithGroups(val)
	if err != nil {
		newCmdResp := protos.CmdResp{
			Id:    cmd.GetId(),
			Name:  cmd.GetName(),
			Error: fmt.Sprintf("Process (%s) does not exist", val),
		}
		return &newCmdResp
	}
	newCmdResp := protos.CmdResp{
		Id:      cmd.GetId(),
		Name:    cmd.GetName(),
		Process: p.ToProtobuf(),
	}
	return &newCmdResp
}
