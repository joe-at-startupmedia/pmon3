package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pmon3/cli/cmd/base"
	"pmon3/cli/cmd/exec"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/test/e2e/cli_helper"
	"strings"
	"testing"

	"pmon3/pmond/protos"

	"time"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3CoreTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3CoreTestSuite))
}

func (suite *Pmon3CoreTestSuite) SetupSuite() {
	suite.cliHelper = cli_helper.SetupSuite(&suite.Suite, "/test/e2e/config/test-config.core.yml", "/test/e2e/config/process.core-test.config.json", "core")
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3CoreTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssert(2)
}

func (suite *Pmon3CoreTestSuite) TestB_AddingAdditionalProcessesFromProcessConfig() {
	ef := model.ExecFlags{
		Name: "test-server-3",
	}
	exec.Exec(suite.cliHelper.ProjectPath+"/test/app/bin/test_app", ef)
	time.Sleep(2 * time.Second)
	passing, _ := suite.cliHelper.LsAssertStatus(3, "running", 0)
	if !passing {
		return
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-4\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(4, "running", 0)

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-5\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(5, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestC1_DescribingAProcessWithAFourthId() {
	newCmdResp := suite.cliHelper.ExecBase1("desc", "4")
	assert.Equal(suite.T(), "test-server-4", newCmdResp.GetProcess().GetName())
}

func (suite *Pmon3CoreTestSuite) TestC2_DescribingANonExistentProcess() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("desc", "6")
	assert.Equal(suite.T(), "process (6) does not exist", newCmdResp.GetError())
}

func (suite *Pmon3CoreTestSuite) TestD1_DeletingAProcess() {
	suite.cliHelper.ExecBase1("del", "3")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(4, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestD2_ForceDeletingAProcess() {
	suite.cliHelper.ExecBase2("del", "4", "force")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(3, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestD3_ForceDeletingANonExistentProcess() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("del", "6")
	time.Sleep(2 * time.Second)
	assert.Equal(suite.T(), "process (6) does not exist", newCmdResp.GetError())
}

func (suite *Pmon3CoreTestSuite) TestE_KillProcesses() {
	var sent *protos.Cmd
	sent = base.SendCmd("kill", "")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		time.Sleep(2 * time.Second)
		suite.cliHelper.LsAssertStatus(3, "stopped", 0)
	}
}

func (suite *Pmon3CoreTestSuite) TestF1_InitAll() {
	var sent *protos.Cmd
	sent = base.SendCmdArg2("init", "all", "blocking")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		time.Sleep(2 * time.Second)
		suite.cliHelper.LsAssertStatus(3, "running", 0)
	}
}

func (suite *Pmon3CoreTestSuite) TestF2_Top() {
	cmdResp := suite.cliHelper.ExecBase0("top")
	pidCsv := cmdResp.GetValueStr()
	assert.Greater(suite.T(), len(pidCsv), 5)

	pids := strings.Split(pidCsv, ",")
	assert.Equal(suite.T(), 4, len(pids))
}

func (suite *Pmon3CoreTestSuite) TestG1_Drop() {
	suite.cliHelper.ExecBase0("drop")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssert(0)
}

func (suite *Pmon3CoreTestSuite) TestG2_DropAfterDrop() {
	suite.cliHelper.ShouldError().ExecBase0("drop")
}

func (suite *Pmon3CoreTestSuite) TestH_InitAllAfterDrop() {
	suite.cliHelper.ExecBase2("init", "all", "blocking")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(2, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestI_StartingAndStopping() {
	suite.cliHelper.ExecBase1("drop", "force")
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-6\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "running", 0)
	suite.cliHelper.ExecBase1("stop", "1")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "stopped", 0)
	suite.cliHelper.ExecBase2("restart", "1", "{}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestJ_RestartIncrementsCounter() {
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-7\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(2, "running", 0)
	suite.cliHelper.ExecBase2("restart", "1", "{}")
	time.Sleep(2 * time.Second)
	newCmdResp := suite.cliHelper.ExecBase1("desc", "1")
	assert.Equal(suite.T(), uint32(2), newCmdResp.GetProcess().GetRestartCount())
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	time.Sleep(2 * time.Second)
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase2("restart", "2", "{}")
	time.Sleep(2 * time.Second)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	assert.Equal(suite.T(), uint32(1), newCmdResp.GetProcess().GetRestartCount())
}

func (suite *Pmon3CoreTestSuite) TestK_ResetRestartCounter() {
	suite.cliHelper.ExecBase1("reset", "1")
	time.Sleep(2 * time.Second)
	newCmdResp := suite.cliHelper.ExecBase1("desc", "1")
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase0("reset")
	time.Sleep(2 * time.Second)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase0("drop")
	time.Sleep(2 * time.Second)
}

func (suite *Pmon3CoreTestSuite) TestL_ResetNonExistentProcess() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("reset", "1")
	assert.Equal(suite.T(), "process (1) does not exist", newCmdResp.GetError())
}

func (suite *Pmon3CoreTestSuite) TestM_ExecProcessWithNonExistentAbsolutePath() {
	newCmdResp := suite.cliHelper.ShouldError().ExecCmd("/nonexistent_path/test_app", "{\"name\": \"test-server-7\"}")
	assert.Contains(suite.T(), newCmdResp.GetError(), "does not exist: stat")
}

func (suite *Pmon3CoreTestSuite) TestN_ExecProcessWithNonExistentRelativePath() {
	ef := model.ExecFlags{}
	execFlags, _ := ef.Parse("{\"name\": \"test-server-7\"}")
	execFlags.File = "./nonexistent_path/test_app"
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("exec", execFlags.Json())
	assert.Contains(suite.T(), newCmdResp.GetError(), "does not exist: stat")
}

func (suite *Pmon3CoreTestSuite) TestO_ExecProcessWithExistentRelativePath() {
	ef := model.ExecFlags{}
	execFlags, _ := ef.Parse("{\"name\": \"test-server-7\"}")
	execFlags.File = "../app/bin/test_app"
	suite.cliHelper.ExecBase1("exec", execFlags.Json())
}

func (suite *Pmon3CoreTestSuite) TestP_ExecProcessWithMalformedJson() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("exec", "{\"malformed\": \"json")
	assert.Contains(suite.T(), newCmdResp.GetError(), "could not parse flags: unexpected end of JSON input")
}

func (suite *Pmon3CoreTestSuite) TestQ_ExecProcessWithoutName() {
	ef := model.ExecFlags{
		File:          suite.cliHelper.ProjectPath + "/test/app/bin/test_app",
		NoAutoRestart: true,
	}
	suite.cliHelper.ExecBase1("exec", ef.Json())
	time.Sleep(2 * time.Second)
}

func (suite *Pmon3CoreTestSuite) TestR_KilledProcessShouldRestart() {

	_, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	p := model.ProcessFromProtobuf(processList[0])

	assert.Greater(suite.T(), len(p.GetPidStr()), 2)

	err := process.SendOsKillSignal(p, true)

	if err != nil {
		suite.Fail(err.Error())
	}

	time.Sleep(6 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)

	if !passing {
		return
	}

	processList = cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), uint32(1), processList[0].GetRestartCount())
}

func (suite *Pmon3CoreTestSuite) TestS_KilledProcessShouldNotRestart() {

	_, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	p := model.ProcessFromProtobuf(processList[1])

	assert.Greater(suite.T(), len(p.GetPidStr()), 2)

	err := process.SendOsKillSignal(p, true)

	if err != nil {
		suite.Fail(err.Error())
	}

	time.Sleep(6 * time.Second)

	_, cmdResp = suite.cliHelper.LsAssertStatus(1, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestT_StartOneFromProcessConfig() {
	suite.cliHelper.ExecBase0("drop")
	time.Sleep(2 * time.Second)
	suite.cliHelper.ExecBase2("restart", "test-server-2", "{}")
	time.Sleep(2 * time.Second)
	_, cmdResp := suite.cliHelper.LsAssertStatus(1, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), 1, len(processList))

	p := model.ProcessFromProtobuf(processList[0])

	assert.Equal(suite.T(), "test-server-2", p.Name)
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3CoreTestSuite) TestZ_TearDown() {
	suite.cliHelper.DropAndClose()
}
