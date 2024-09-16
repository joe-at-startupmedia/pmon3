package group

import (
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"strings"
)

func Assign(cmd *protos.Cmd) *protos.CmdResp {

	groupNameOrId := strings.Split(cmd.GetArg1(), ",")
	processNameOrId := strings.Split(cmd.GetArg2(), ",")

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}

	var groups []*model.Group
	var processes []*model.Process

	for i := range groupNameOrId {
		group, err := repo.Group().FindByIdOrName(groupNameOrId[i])
		if err != nil {
			newCmdResp.Error = err.Error()
			return &newCmdResp
		}
		groups = append(groups, group)
	}

	for i := range processNameOrId {
		process, err := repo.Process().FindByIdOrName(processNameOrId[i])
		if err != nil {
			newCmdResp.Error = err.Error()
			return &newCmdResp
		}
		processes = append(processes, process)
	}

	for i := range groups {
		for j := range processes {
			err := repo.GroupOf(groups[i]).AssignProcess(processes[j])
			if err != nil {
				newCmdResp.Error = err.Error()
				return &newCmdResp
			}
		}
	}

	return &newCmdResp
}
