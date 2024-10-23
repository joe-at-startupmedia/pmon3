package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	"pmon3/protos"
	"strings"
)

func Create(cmd *protos.Cmd) *protos.CmdResp {
	groupName := cmd.GetArg1()
	_, err := repo.Group().Create(groupName)
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	if err != nil {
		if strings.HasPrefix(err.Error(), "command error:") {
			return base.ErroredCmdResp(cmd, err)
		} else {
			newCmdResp.Error = err.Error()
		}
	}
	return &newCmdResp
}
