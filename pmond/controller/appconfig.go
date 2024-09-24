package controller

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func AppConfig(cmd *protos.Cmd) *protos.CmdResp {
	all, err := repo.Process().FindAll()
	if err != nil {
		pmond.Log.Fatalf("pmon3 can find processes: %v", err)
	}
	appsConfig := model.AppsConfig{
		Apps: make([]model.AppsConfigApp, len(all)),
	}
	for i, p := range all {
		appsConfig.Apps[i] = *p.ToAppsConfigApp()
	}

	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: appsConfig.Json(),
	}
	return &newCmdResp
}
