package controller

import (
	"fmt"
	"github.com/joe-at-startupmedia/pmq_responder"
	"google.golang.org/protobuf/proto"
	"pmon3/pmond/protos"
)

func ErroredCmdResp(cmd *protos.Cmd, err error) *protos.CmdResp {
	return &protos.CmdResp{
		Id:    cmd.GetId(),
		Name:  cmd.GetName(),
		Error: fmt.Sprintf("%s, %s", cmd.GetName(), err),
	}
}

// HandleCmdRequest provides a concrete implementation of HandleRequestFromProto using the local Cmd protobuf type
func HandleCmdRequest(mqr *pmq_responder.MqResponder) error {

	cmd := &protos.Cmd{}
	return mqr.HandleRequestFromProto(cmd, func() (processed []byte, err error) {

		var cmdResp *protos.CmdResp
		switch cmd.GetName() {
		case "log":
		case "logf":
		case "desc":
			cmdResp = Desc(cmd)
		case "list":
			cmdResp = List(cmd)
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
	})
}
