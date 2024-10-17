package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/restart"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Restart(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()
	flags := cmd.GetArg2()
	return RestartByParams(cmd, idOrName, flags, true)
}

func RestartByParams(cmd *protos2.Cmd, idOrName string, flags string, incrementCounter bool) *protos2.CmdResp {
	// kill the process and insert a new record with "queued" status

	p, err := repo.Process().FindByIdOrName(idOrName)

	//the process doesn't exist,  so we'll look in the AppConfig
	if err != nil {

		p, err = restart.ByProcess(cmd, nil, idOrName, flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}

		newCmdResp := protos2.CmdResp{
			Id:      cmd.GetId(),
			Name:    cmd.GetName(),
			Process: p.ToProtobuf(),
		}
		return &newCmdResp
	} else {
		_, err = restart.ByProcess(cmd, p, idOrName, flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		} else {
			newCmdResp := protos2.CmdResp{
				Id:      cmd.GetId(),
				Name:    cmd.GetName(),
				Process: p.ToProtobuf(),
			}
			return &newCmdResp
		}
	}
}
