package base

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/joe-at-startupmedia/xipc"
	"os/user"
	"pmon3/cli"
	"pmon3/pmond/protos"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var xr xipc.IRequester

const SEND_RECEIVE_TIMEOUT = time.Second * 5

var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Margin(1)

func handleOpenError(e error) {
	if e != nil {
		if strings.Contains(e.Error(), "Could not apply permissions") {
			cli.Log.Debugf("could not apply sender permissions: %s", e.Error())
		} else {
			cli.Log.Fatal("could not initialize sender: ", e.Error())
		}
	}
}

func sendCmd(cmd *protos.Cmd) {
	cli.Log.Debug("sending message")
	pbm := proto.Message(cmd)
	sendErrChan := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), SEND_RECEIVE_TIMEOUT)
	defer cancel()

	go func() {
		err := xr.RequestUsingProto(&pbm)
		sendErrChan <- err
	}()

	select {
	case <-ctx.Done():
		cli.Log.Fatal("operation timed out")
	case err := <-sendErrChan:
		if err != nil {
			cli.Log.Fatal(err)
		}
		cli.Log.Debugf("Sent a new message: %s", cmd.String())
	}
}

func GetResponse(sent *protos.Cmd) *protos.CmdResp {
	cli.Log.Debug("getting response")
	newCmdResp := &protos.CmdResp{}
	readErrChan := make(chan error, 1)

	sendRcvTimeout := SEND_RECEIVE_TIMEOUT
	confExecResponseWait := cli.Config.GetCmdExecResponseWait()

	//because processes with dependencies take longer to start
	if sent.GetName() == "init" {
		sendRcvTimeout = time.Second * 180
	} else if confExecResponseWait > sendRcvTimeout {
		sendRcvTimeout = confExecResponseWait
	}

	ctx, cancel := context.WithTimeout(context.Background(), sendRcvTimeout)
	defer cancel()

	go func() {
		timer := time.NewTicker(time.Millisecond * 100)
		for {
			_, err := waitForResponse(newCmdResp)
			if newCmdResp.GetId() != sent.GetId() {
				cli.Log.Debugf("response (%s) doesn't match sent (%s). skipping.", newCmdResp.GetId(), sent.GetId())
				<-timer.C
				continue
			}
			readErrChan <- err
			break
		}

	}()

	select {
	case <-ctx.Done():
		cli.Log.Fatal("operation timed out")
	case err := <-readErrChan:
		if err != nil {
			cli.Log.Fatal(err)
		} else if len(newCmdResp.GetError()) > 0 {
			OutputError(newCmdResp.GetError())
		}
		cli.Log.Debugf("Got a new response: %s", newCmdResp.Id)
	}
	return newCmdResp
}

func OutputError(error string) {
	fmt.Println(errorStyle.Render(error))
}

func SendCmd(cmdName string, arg1 string) *protos.Cmd {
	cmd := &protos.Cmd{
		Id:   uuid.NewString(),
		Name: cmdName,
		Arg1: arg1,
	}

	sendCmd(cmd)

	return cmd
}

func SendCmdArg2(cmdName string, arg1 string, arg2 string) *protos.Cmd {
	cmd := &protos.Cmd{
		Id:   uuid.NewString(),
		Name: cmdName,
		Arg1: arg1,
		Arg2: arg2,
	}

	sendCmd(cmd)

	return cmd
}

func waitForResponse(newCmdResp *protos.CmdResp) (*proto.Message, error) {
	return xr.WaitForProto(newCmdResp)
}

func CloseSender() error {
	return xr.CloseRequester()
}

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		cli.Log.Errorf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}
