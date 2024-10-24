package group

import (
	"fmt"
	"pmon3/model"
	"pmon3/pmond/repo"
	"pmon3/protos"
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
	return listProcesses(cmd, g)
}

func listProcesses(cmd *protos.Cmd, g *model.Group) *protos.CmdResp {
	newProcessList := protos.ProcessList{}
	for _, p := range g.Processes {
		p.SetUsageStats()
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
