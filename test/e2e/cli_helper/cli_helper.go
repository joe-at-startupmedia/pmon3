package cli_helper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"strings"

	"time"
)

type CliHelper struct {
	suite       *suite.Suite
	ProjectPath string
	shouldError bool
}

func New(suite *suite.Suite, projectPath string) *CliHelper {
	return &CliHelper{
		suite,
		projectPath,
		false,
	}
}

func (cliHelper *CliHelper) reset() {
	cliHelper.shouldError = false
}

func (cliHelper *CliHelper) ShouldError() *CliHelper {
	cliHelper.shouldError = true
	return cliHelper
}

func (cliHelper *CliHelper) LsAssert(expectedProcessLen int) (bool, *protos.CmdResp) {
	newCmdResp := cliHelper.ExecBase0("list")
	processList := newCmdResp.GetProcessList().GetProcesses()
	cli.Log.Infof("process list: %s \n value string: %s \n", processList, newCmdResp.GetValueStr())
	passing := assert.Equal(cliHelper.suite.T(), expectedProcessLen, len(processList))
	//cli.Log.Fatalf("Expected process length of %d but got %d", expectedProcessLen, actualProcessLen)
	return passing, newCmdResp
}

func (cliHelper *CliHelper) LsAssertStatus(expectedProcessLen int, status string, retries int) (bool, *protos.CmdResp) {

	if retries == 3 {
		cliHelper.suite.Fail("assert status failed with maximum of 3 retries")
	}

	cmdResp := cliHelper.ExecBase0("list")
	processList := cmdResp.GetProcessList().GetProcesses()

	processesMatchingStatus := 0

	for _, p := range processList {
		if p.Status == status {
			processesMatchingStatus++
		}
	}

	if expectedProcessLen != processesMatchingStatus && retries < 3 {
		cli.Log.Warnf("retry count: %d with params: %d %d %s", retries+1, expectedProcessLen, processesMatchingStatus, status)
		time.Sleep(time.Second * 5)
		return cliHelper.LsAssertStatus(expectedProcessLen, status, retries+1)
	}

	passing := assert.Equal(cliHelper.suite.T(), expectedProcessLen, processesMatchingStatus)

	return passing, cmdResp
}

func (cliHelper *CliHelper) ExecCmd(processFile string, execFlagsJson string) *protos.CmdResp {
	processFile = cliHelper.ProjectPath + processFile
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
	if len(newCmdResp.GetError()) > 0 && !cliHelper.shouldError {
		cliHelper.suite.Fail(newCmdResp.GetError())
	}
	cliHelper.reset()
	return newCmdResp
}

func (cliHelper *CliHelper) DgraphProcessNames(arg1 string) ([]string, []string) {
	cmdResp := cliHelper.ExecBase1("dgraph", arg1)

	response := strings.Split(cmdResp.GetValueStr(), "||")

	var nonDeptProcessNames []string
	var deptProcessNames []string
	if len(response[0]) > 0 {
		nonDeptProcessNames = strings.Split(response[0], "\n")
	}
	if len(response[1]) > 0 {
		deptProcessNames = strings.Split(response[1], "\n")
	}

	return nonDeptProcessNames, deptProcessNames
}

func (cliHelper *CliHelper) ShouldKill(expectedProcessLen int, waitBeforeAssertion int) bool {
	cliHelper.ExecBase0("kill")
	time.Sleep(time.Duration(waitBeforeAssertion) * time.Second)
	passing, _ := cliHelper.LsAssertStatus(expectedProcessLen, "stopped", 0)
	return passing
}
