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

func HandleMessage(msg []byte) ([]byte, error) {
	newCmd := &protos.Cmd{}
	if err := proto.Unmarshal(msg, newCmd); err != nil {
		pmond.Log.Fatal("unmarshaling error: ", err)
	}

	pmond.Log.Debug(newCmd.String())

	var newCmdResp *protos.CmdResp

	switch newCmd.GetName() {
	case "log":
	case "logf":
	case "desc":
		newCmdResp = Desc(newCmd)
	case "list":
		newCmdResp = List(newCmd)
	case "restart":
		newCmdResp = Restart(newCmd)
	case "exec":
		newCmdResp = Exec(newCmd)
	case "stop":
		newCmdResp = Stop(newCmd)
	case "kill":
		newCmdResp = Kill(newCmd)
	case "init":
		newCmdResp = Initialize(newCmd)
	case "del":
		newCmdResp = Delete(newCmd)
	case "drop":
		newCmdResp = Drop(newCmd)
	}

	processed, err := proto.Marshal(newCmdResp)

	if err != nil {
		pmond.Log.Fatal("marshaling error: ", err)
	}

	return processed, err
}
