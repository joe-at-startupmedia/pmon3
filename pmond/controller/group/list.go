package group

import (
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func List(cmd *protos.Cmd) *protos.CmdResp {
	all, err := repo.Group().FindAll()
	if err != nil {
		pmond.Log.Fatalf("pmon3 can find processes: %v", err)
	}
	newGroupList := protos.GroupList{}
	for _, g := range all {
		newGroupList.Groups = append(newGroupList.Groups, g.ToProtobuf())
	}
	newCmdResp := protos.CmdResp{
		Id:        cmd.GetId(),
		Name:      cmd.GetName(),
		GroupList: &newGroupList,
	}
	return &newCmdResp
}
