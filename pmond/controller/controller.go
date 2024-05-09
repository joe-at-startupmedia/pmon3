package controller

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"pmon3/pmond"
	"pmon3/pmond/protos"
)

func ErroredCmdResp(cmd *protos.Cmd, err error) *protos.CmdResp {
	return &protos.CmdResp{
		Id:    cmd.GetId(),
		Name:  cmd.GetName(),
		Error: fmt.Sprintf("%s, %s", cmd.GetName(), err),
	}
}

func MsgHandler(cmd *protos.Cmd) (processed []byte, err error) {

	pmond.Log.Infof("got a cmd: %s", cmd)
	var cmdResp *protos.CmdResp
	switch cmd.GetName() {
	case "log":
		fallthrough
	case "logf":
		fallthrough
	case "desc":
		cmdResp = Desc(cmd)
	case "list":
		cmdResp = List(cmd)
	case "top":
		cmdResp = Top(cmd)
	case "restart":
		cmdResp = Restart(cmd)
	case "exec":
		cmdResp = Exec(cmd)
	case "stop":
		cmdResp = Stop(cmd)
	case "kill":
		cmdResp = Kill(cmd)
	case "init":
		cmdResp = Initialize(cmd)
	case "del":
		cmdResp = Delete(cmd)
	case "drop":
		cmdResp = Drop(cmd)
	}

	data, err := proto.Marshal(cmdResp)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: %w", err)
	}

	return data, nil
}
