package cli_helper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"

	"time"
)

type CliHelper struct {
	suite      *suite.Suite
	appBinPath string
}

func New(suite *suite.Suite, appBinPath string) *CliHelper {
	return &CliHelper{
		suite,
		appBinPath,
	}
}

func (cliHelper *CliHelper) LsAssert(expectedProcessLen int) *protos.CmdResp {
	newCmdResp := cliHelper.ExecBase0("list")
	processList := newCmdResp.GetProcessList().GetProcesses()
	cli.Log.Infof("process list: %s \n value string: %s \n", processList, newCmdResp.GetValueStr())
	assert.Equal(cliHelper.suite.T(), expectedProcessLen, len(processList))
	//cli.Log.Fatalf("Expected process length of %d but got %d", expectedProcessLen, actualProcessLen)
	return newCmdResp
}

func (cliHelper *CliHelper) LsAssertStatus(expectedProcessLen int, status string, retries int) {

	cmdResp := cliHelper.LsAssert(expectedProcessLen)
	processList := cmdResp.GetProcessList().GetProcesses()

	for _, p := range processList {
		if p.Status != status && retries < 3 { //three retries are allowed
			cli.Log.Infof("Expected process status of %s but got %s", status, p.Status)
			cli.Log.Warnf("retry count: %d", retries+1)
			time.Sleep(time.Second * 5)
			cliHelper.LsAssertStatus(expectedProcessLen, status, retries+1)
			break
		} else {
			assert.Equal(cliHelper.suite.T(), status, p.Status)
		}
	}
}

func (cliHelper *CliHelper) ExecCmd(processFile string, execFlagsJson string) *protos.CmdResp {
	processFile = cliHelper.appBinPath + processFile
	cli.Log.Infof("Executing: pmon3 exec %s %s", processFile, execFlagsJson)
	ef := model.ExecFlags{}
	execFlags, err := ef.Parse(execFlagsJson)
	if err != nil {
		cli.Log.Fatal(err)
	}
	execFlags.File = processFile
	return cliHelper.ExecBase1("exec", execFlags.Json())
}

func (cliHelper *CliHelper) ExecBase0(cmd string) *protos.CmdResp {
	return cliHelper.execBase(cmd, "", "")
}

func (cliHelper *CliHelper) ExecBase1(cmd string, arg1 string) *protos.CmdResp {
	return cliHelper.execBase(cmd, arg1, "")
}

func (cliHelper *CliHelper) ExecBase2(cmd string, arg1 string, arg2 string) *protos.CmdResp {
	return cliHelper.execBase(cmd, arg1, arg2)
}

func (cliHelper *CliHelper) execBase(cmd string, arg1 string, arg2 string) *protos.CmdResp {
	var sent *protos.Cmd
	if len(arg2) > 0 {
		sent = base.SendCmdArg2(cmd, arg1, arg2)
	} else {
		sent = base.SendCmd(cmd, arg1)
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cliHelper.suite.Fail(newCmdResp.GetError())
	}
	return newCmdResp
}
