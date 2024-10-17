package e2e

import (
	"context"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli/controller"
	model2 "pmon3/model"
	"pmon3/pmond/process"
	"pmon3/test/e2e/cli_helper"
	"strings"
	"sync"
	"testing"

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

func (suite *Pmon3CoreTestSuite) Sleep() {
	time.Sleep(suite.cliHelper.GetSleepDurationFromEnv(0, "core"))
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3CoreTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssert(2)
}

func (suite *Pmon3CoreTestSuite) TestB_AddingAdditionalProcessesFromProcessConfig() {
	ef := model2.ExecFlags{
		Name: "test-server-3",
	}
	controller.Exec(suite.cliHelper.ProjectPath+"/test/app/bin/test_app", ef)
	suite.Sleep()
	passing, _ := suite.cliHelper.LsAssertStatus(3, "running", 0)
	if !passing {
		return
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-4\"}")
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(4, "running", 0)

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-5\"}")
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(5, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestC1_DescribingAProcessWithAFourthId() {
	newCmdResp := controller.Desc("4")
	assert.Equal(suite.T(), "test-server-4", newCmdResp.GetProcess().GetName())
}

func (suite *Pmon3CoreTestSuite) TestC2_DescribingANonExistentProcess() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("desc", "6")
	assert.Equal(suite.T(), "process (6) does not exist", newCmdResp.GetError())
}

func (suite *Pmon3CoreTestSuite) TestD1_DeletingAProcess() {
	controller.Del("3", false)
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(4, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestD2_ForceDeletingAProcess() {
	controller.Del("4", true)
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(3, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestD3_ForceDeletingANonExistentProcess() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("del", "6")
	suite.Sleep()
	assert.Equal(suite.T(), "process (6) does not exist", newCmdResp.GetError())
}

func (suite *Pmon3CoreTestSuite) TestE_KillProcesses() {
	newCmdResp := controller.Kill(false)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		suite.Sleep()
		suite.cliHelper.LsAssertStatus(3, "stopped", 0)
	}
}

func (suite *Pmon3CoreTestSuite) TestF1_InitAll() {
	newCmdResp := controller.Initialize(false, true)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		suite.Sleep()
		suite.cliHelper.LsAssertStatus(3, "running", 0)
	}
}

func onKeyboardEventSort() chan controller.KeyboardResult {
	ch := make(chan controller.KeyboardResult)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- controller.KeyboardResult{
			Char: 's',
		}
	}()
	return ch
}

func onKeyboardEventEscape() chan controller.KeyboardResult {
	ch := make(chan controller.KeyboardResult)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- controller.KeyboardResult{
			Key: keyboard.KeyEsc,
		}
	}()
	return ch
}

func onKeyboardEventError() chan controller.KeyboardResult {
	ch := make(chan controller.KeyboardResult)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- controller.KeyboardResult{
			Err: fmt.Errorf("simulating an error for testing"),
		}
	}()
	return ch
}

func (suite *Pmon3CoreTestSuite) TestF2_Top() {

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go controller.Topn(2, ctx, &wg, onKeyboardEventSort, os.Stdout)
	suite.cliHelper.SleepFor(time.Millisecond * 6000)
	cancel() //will call wg.Done

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go controller.Topn(2, context.Background(), &wg2, onKeyboardEventEscape, os.Stdout)
	suite.cliHelper.SleepFor(time.Millisecond * 4000)
	//cancel() //will call wg.Done

	var wg3 sync.WaitGroup
	wg3.Add(1)
	go controller.Topn(2, context.Background(), &wg3, onKeyboardEventError, os.Stdout)
	suite.cliHelper.SleepFor(time.Millisecond * 4000)

	cmdResp := suite.cliHelper.ExecBase0("top")
	pidCsv := cmdResp.GetValueStr()
	assert.Greater(suite.T(), len(pidCsv), 5)

	pids := strings.Split(pidCsv, ",")
	assert.Equal(suite.T(), 4, len(pids))
}

func (suite *Pmon3CoreTestSuite) TestG1_Drop() {
	controller.Drop(false)
	suite.Sleep()
	suite.cliHelper.SleepFor(time.Millisecond * 1000)
	suite.cliHelper.LsAssert(0)
}

func (suite *Pmon3CoreTestSuite) TestG2_DropAfterDrop() {
	suite.cliHelper.ShouldError().ExecBase0("drop")
}

func (suite *Pmon3CoreTestSuite) TestH_InitAllAfterDrop() {
	suite.cliHelper.ExecBase2("init", "all", "blocking")
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(2, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestI_StartingAndStopping() {
	controller.Drop(true)
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-6\"}")
	suite.Sleep()
	controller.List()
	suite.cliHelper.LsAssertStatus(1, "running", 0)
	controller.Stop("1", false)
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(1, "stopped", 0)
	controller.Restart("restart", "1", "{}")
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(1, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestJ_RestartIncrementsCounter() {
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-7\"}")
	suite.Sleep()
	suite.cliHelper.LsAssertStatus(2, "running", 0)
	suite.cliHelper.ExecBase2("restart", "1", "{}")
	suite.Sleep()
	newCmdResp := suite.cliHelper.ExecBase1("desc", "1")
	assert.Equal(suite.T(), uint32(2), newCmdResp.GetProcess().GetRestartCount())
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	suite.Sleep()
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase2("restart", "2", "{}")
	suite.Sleep()
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	assert.Equal(suite.T(), uint32(1), newCmdResp.GetProcess().GetRestartCount())
}

func (suite *Pmon3CoreTestSuite) TestK_ResetRestartCounter() {
	controller.Reset("1")
	suite.Sleep()
	newCmdResp := suite.cliHelper.ExecBase1("desc", "1")
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase0("reset")
	suite.Sleep()
	newCmdResp = suite.cliHelper.ExecBase1("desc", "2")
	assert.Equal(suite.T(), uint32(0), newCmdResp.GetProcess().GetRestartCount())
	suite.cliHelper.ExecBase0("drop")
	suite.Sleep()
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
	ef := model2.ExecFlags{}
	execFlags, _ := ef.Parse("{\"name\": \"test-server-7\"}")
	execFlags.File = "./nonexistent_path/test_app"
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("exec", execFlags.Json())
	assert.Contains(suite.T(), newCmdResp.GetError(), "does not exist: stat")
}

func (suite *Pmon3CoreTestSuite) TestO_ExecProcessWithExistentRelativePath() {
	ef := model2.ExecFlags{}
	execFlags, _ := ef.Parse("{\"name\": \"test-server-7\"}")
	execFlags.File = "../app/bin/test_app"
	suite.cliHelper.ExecBase1("exec", execFlags.Json())
}

func (suite *Pmon3CoreTestSuite) TestP_ExecProcessWithMalformedJson() {
	newCmdResp := suite.cliHelper.ShouldError().ExecBase1("exec", "{\"malformed\": \"json")
	assert.Contains(suite.T(), newCmdResp.GetError(), "could not parse flags: unexpected end of JSON input")
}

func (suite *Pmon3CoreTestSuite) TestQ_ExecProcessWithoutName() {
	ef := model2.ExecFlags{
		File:          suite.cliHelper.ProjectPath + "/test/app/bin/test_app",
		NoAutoRestart: true,
	}
	suite.cliHelper.ExecBase1("exec", ef.Json())
	suite.Sleep()
}

func (suite *Pmon3CoreTestSuite) TestR_KilledProcessShouldRestart() {

	_, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	p := model2.ProcessFromProtobuf(processList[0])

	assert.Greater(suite.T(), len(p.GetPidStr()), 2)

	err := process.SendOsKillSignal(p, true)

	if err != nil {
		suite.Fail(err.Error())
	}

	suite.cliHelper.SleepFor(time.Second * 5)

	expectedRestartCount := uint32(1)
	var actualProcessCount uint32

	for range 3 {
		passing, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)
		if !passing {
			break
		}
		processList = cmdResp.GetProcessList().GetProcesses()
		actualProcessCount = processList[0].GetRestartCount()
		if expectedRestartCount == actualProcessCount {
			break
		}
		suite.cliHelper.SleepFor(time.Second * 1)
	}
	assert.Equal(suite.T(), expectedRestartCount, actualProcessCount)
}

func (suite *Pmon3CoreTestSuite) TestS_KilledProcessShouldNotRestart() {

	_, cmdResp := suite.cliHelper.LsAssertStatus(2, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	p := model2.ProcessFromProtobuf(processList[1])

	assert.Greater(suite.T(), len(p.GetPidStr()), 2)

	err := process.SendOsKillSignal(p, true)

	if err != nil {
		suite.Fail(err.Error())
	}

	suite.Sleep()
	//suite.cliHelper.SleepFor(time.Second * 5)

	_, cmdResp = suite.cliHelper.LsAssertStatus(1, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestT_StartOneFromProcessConfig() {
	suite.cliHelper.ExecBase0("drop")
	suite.Sleep()
	suite.cliHelper.ExecBase2("restart", "test-server-2", "{}")
	suite.Sleep()
	_, cmdResp := suite.cliHelper.LsAssertStatus(1, "running", 0)

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), 1, len(processList))

	p := model2.ProcessFromProtobuf(processList[0])

	assert.Equal(suite.T(), "test-server-2", p.Name)
}

func (suite *Pmon3CoreTestSuite) TestU_LogProcess() {
	cmdResp := controller.Log("test-server-2", true, "10")
	assert.Equal(suite.T(), len(cmdResp.GetError()), 0)
}

func (suite *Pmon3CoreTestSuite) TestW_LogfProcess() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cmdResp := controller.Logf("test-server-2", "10", ctx)
		assert.Equal(suite.T(), len(cmdResp.GetError()), 0)
	}()
	suite.cliHelper.SleepFor(time.Millisecond * 1000)
	cancel()
}

func (suite *Pmon3CoreTestSuite) TestW_Stop() {
	controller.Stop("1", true)
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3CoreTestSuite) TestZ_TearDown() {
	suite.cliHelper.DropAndClose()
}
