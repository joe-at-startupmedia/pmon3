package group

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
	g, err := repo.Group().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	for i := range g.Processes {
		p := g.Processes[i]
		err = restart.ByProcess(cmd, p, idOrName, flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}
	}
	return listProcesses(cmd, g)
}
