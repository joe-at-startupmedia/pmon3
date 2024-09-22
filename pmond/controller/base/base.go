package base

import (
	"fmt"
	"pmon3/pmond/protos"
)

func ErroredCmdResp(cmd *protos.Cmd, err error) *protos.CmdResp {
	return &protos.CmdResp{
		Id:    cmd.GetId(),
		Name:  cmd.GetName(),
		Error: fmt.Sprintf("%s, %s", cmd.GetName(), err),
	}
}
