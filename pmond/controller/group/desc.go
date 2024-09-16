package group

import (
	"fmt"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Desc(cmd *protos.Cmd) *protos.CmdResp {
	val := cmd.GetArg1()
	g, err := repo.Group().FindByIdOrName(val)
	if err != nil {
		newCmdResp := protos.CmdResp{
			Id:    cmd.GetId(),
			Name:  cmd.GetName(),
			Error: fmt.Sprintf("Group (%s) does not exist", val),
		}
		return &newCmdResp
	}
	newProcessList := protos.ProcessList{}
	for _, p := range g.Processes {
		newProcessList.Processes = append(newProcessList.Processes, p.ToProtobuf())
	}
	newCmdResp := protos.CmdResp{
		Id:          cmd.GetId(),
		Name:        cmd.GetName(),
		Group:       g.ToProtobuf(),
		ProcessList: &newProcessList,
	}
	return &newCmdResp
}
