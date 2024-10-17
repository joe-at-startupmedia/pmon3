package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func List(cmd *protos2.Cmd) *protos2.CmdResp {
	all, err := repo.Group().FindAll()
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}
	newGroupList := protos2.GroupList{}
	for _, g := range all {
		newGroupList.Groups = append(newGroupList.Groups, g.ToProtobuf())
	}
	newCmdResp := protos2.CmdResp{
		Id:        cmd.GetId(),
		Name:      cmd.GetName(),
		GroupList: &newGroupList,
	}
	return &newCmdResp
}
