package controller

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Desc(cmd *protos.Cmd) *protos.CmdResp {
	val := cmd.GetArg1()
	err, p := model.FindProcessByIdOrName(pmond.Db(), val)
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
