package cli_helper

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli"
	"pmon3/cli/controller/base"
	"pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/god"
	"pmon3/protos"
	"pmon3/utils/conv"
	"strings"

	"time"
)

type CliHelper struct {
	suite        *suite.Suite
	cancelGodCtx context.CancelFunc
	ProjectPath  string
	ArtifactPath string
	shouldError  bool
}

func SetupSuite(s *suite.Suite, configFile string, processConfigFile string, messageQueueSuffix string) *CliHelper {
	projectPath := os.Getenv("PROJECT_PATH")
	artifactPath := os.Getenv("ARTIFACT_PATH")

	if err := pmond.Instance(projectPath+configFile, projectPath+processConfigFile); err != nil {
		s.FailNow(err.Error())
	}
	pmond.Config.MessageQueue.NameSuffix = messageQueueSuffix

	ctx, cancel := context.WithCancel(context.Background())
	go god.Summon(ctx)

	time.Sleep(5 * time.Second)

	if err := cli.Instance(projectPath + configFile); err != nil {
		s.FailNow(err.Error())
	}
	cli.Config.MessageQueue.NameSuffix = messageQueueSuffix

	base.OpenSender()

	fmt.Println(pmond.Config.Yaml())

	return New(s, projectPath, artifactPath, cancel)
}

func New(suite *suite.Suite, projectPath string, artifactPath string, cancelGodCtx context.CancelFunc) *CliHelper {
	return &CliHelper{
		suite,
		cancelGodCtx,
		projectPath,
		artifactPath,
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
		time.Sleep(time.Second * 2)
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

func (cliHelper *CliHelper) GetSleepDurationFromEnv(defaultDuration int, suiteName string) time.Duration {
	var sleepDuration = defaultDuration
	sleepEnvArg := os.Getenv("TEST_SLEEP_" + strings.ToUpper(suiteName))
	if sleepEnvArg != "" {
		sleepDuration = conv.StrToInt(sleepEnvArg)
	} else {
		sleepEnvArg = os.Getenv("TEST_SLEEP")
		if sleepEnvArg != "" {
			sleepDuration = conv.StrToInt(sleepEnvArg)
		}
	}
	return time.Duration(sleepDuration) * time.Millisecond
}

func (cliHelper *CliHelper) SleepFor(duration time.Duration) {
	time.Sleep(duration)
}

func (cliHelper *CliHelper) DropAndClose() {
	cliHelper.ExecBase1("drop", "force")
	cliHelper.Close()
}

func (cliHelper *CliHelper) Close() {
	time.Sleep(3 * time.Second)
	god.Banish()
	cliHelper.cancelGodCtx()
	base.CloseSender()
	time.Sleep(1 * time.Second)
}
