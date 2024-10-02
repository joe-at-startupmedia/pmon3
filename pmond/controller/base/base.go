package base

import (
	"pmon3/pmond/protos"
)

func ErroredCmdResp(cmd *protos.Cmd, err error) *protos.CmdResp {
	return &protos.CmdResp{
		Id:    cmd.GetId(),
		Name:  cmd.GetName(),
		Error: err.Error(),
	}
}
