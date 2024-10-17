package base

import (
	protos2 "pmon3/protos"
)

func ErroredCmdResp(cmd *protos2.Cmd, err error) *protos2.CmdResp {
	return &protos2.CmdResp{
		Id:    cmd.GetId(),
		Name:  cmd.GetName(),
		Error: err.Error(),
	}
}
