package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
	"strings"
)

func Create(cmd *protos2.Cmd) *protos2.CmdResp {
	groupName := cmd.GetArg1()
	_, err := repo.Group().Create(groupName)
	newCmdResp := protos2.CmdResp{
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
