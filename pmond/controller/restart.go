package controller

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/restart"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Restart(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	flags := cmd.GetArg2()
	return RestartByParams(cmd, idOrName, flags, true)
}

func RestartByParams(cmd *protos.Cmd, idOrName string, flags string, incrementCounter bool) *protos.CmdResp {
	// kill the process and insert a new record with "queued" status

	p, err := repo.Process().FindByIdOrName(idOrName)

	//the process doesn't exist,  so we'll look in the AppConfig
	if err != nil {

		err = restart.ByProcess(cmd, nil, idOrName, flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}

		newCmdResp := protos.CmdResp{
			Id:      cmd.GetId(),
			Name:    cmd.GetName(),
			Process: p.ToProtobuf(),
		}
		return &newCmdResp
	} else {
		err = restart.ByProcess(cmd, p, idOrName, flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		} else {
			newCmdResp := protos.CmdResp{
				Id:   cmd.GetId(),
				Name: cmd.GetName(),
				Process: &protos.Process{
					Log: p.Log,
				},
			}
			return &newCmdResp
		}
	}
}
