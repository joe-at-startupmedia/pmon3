package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
	"strings"
)

func Delete(cmd *protos2.Cmd) *protos2.CmdResp {
	groupIdOrName := cmd.GetArg1()
	err := repo.Group().Delete(groupIdOrName)
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
