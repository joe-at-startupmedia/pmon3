package group

import (
	model2 "pmon3/model"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
	"strings"
)

func Remove(cmd *protos2.Cmd) *protos2.CmdResp {

	groupNameOrId := strings.Split(cmd.GetArg1(), ",")
	processNameOrId := strings.Split(cmd.GetArg2(), ",")

	newCmdResp := protos2.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}

	var groups []*model2.Group
	var processes []*model2.Process

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
			err := repo.GroupOf(groups[i]).RemoveProcess(processes[j])
			if err != nil {
				newCmdResp.Error = err.Error()
				return &newCmdResp
			}
		}
	}

	return &newCmdResp
}
