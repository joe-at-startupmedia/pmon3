package group

import (
	"fmt"
	"pmon3/model"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Desc(cmd *protos2.Cmd) *protos2.CmdResp {
	val := cmd.GetArg1()
	g, err := repo.Group().FindByIdOrName(val)
	if err != nil {
		newCmdResp := protos2.CmdResp{
			Id:    cmd.GetId(),
			Name:  cmd.GetName(),
			Error: fmt.Sprintf("Group (%s) does not exist", val),
		}
		return &newCmdResp
	}
	return listProcesses(cmd, g)
}

func listProcesses(cmd *protos2.Cmd, g *model.Group) *protos2.CmdResp {
	newProcessList := protos2.ProcessList{}
	for _, p := range g.Processes {
		newProcessList.Processes = append(newProcessList.Processes, p.ToProtobuf())
	}
	newCmdResp := protos2.CmdResp{
		Id:          cmd.GetId(),
		Name:        cmd.GetName(),
		Group:       g.ToProtobuf(),
		ProcessList: &newProcessList,
	}
	return &newCmdResp
}
