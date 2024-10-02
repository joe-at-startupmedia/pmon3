package e2e

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond"
	"pmon3/pmond/god"
	"pmon3/test/e2e/cli_helper"
	"testing"

	"pmon3/pmond/protos"

	"time"
)

// Define the suite, and absorb the built-in core suite
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

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Pmon3CoreTestSuite) SetupSuite() {

	projectPath := os.Getenv("PROJECT_PATH")
	suite.cliHelper = cli_helper.New(&suite.Suite, projectPath)

	configFile := projectPath + "/test/e2e/config/test-config.yml"
	processConfigFile := projectPath + "/test/e2e/config/process.core-test.config.json"
	if err := cli.Instance(configFile); err != nil {
		suite.FailNow(err.Error())
	}

	if err := pmond.Instance(configFile, processConfigFile); err != nil {
		suite.FailNow(err.Error())
	}

	ctx := context.Background()
	go god.Summon(ctx)

	time.Sleep(5 * time.Second)

	base.OpenSender()
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3CoreTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssert(2)
}

func (suite *Pmon3CoreTestSuite) TestB_AddingAdditionalProcessesFromProcessConfig() {
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-3\"}")
	time.Sleep(2 * time.Second)
	passing, _ := suite.cliHelper.LsAssertStatus(3, "running", 0)
	if !passing {
		return
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-4\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(4, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestC_DescribingAProcessWithAFourthId() {
	newCmdResp := suite.cliHelper.ExecBase1("desc", "4")
	assert.Equal(suite.T(), "test-server-4", newCmdResp.GetProcess().GetName())
}

func (suite *Pmon3CoreTestSuite) TestD_DeletingAProcess() {
	suite.cliHelper.ExecBase1("del", "3")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(3, "running", 0)
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

func (suite *Pmon3CoreTestSuite) TestF_InitAll() {
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

func (suite *Pmon3CoreTestSuite) TestG_Drop() {
	suite.cliHelper.ExecBase0("drop")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssert(0)
}

func (suite *Pmon3CoreTestSuite) TestH_InitAllAfterDrop() {
	suite.cliHelper.ExecBase2("init", "all", "blocking")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(2, "running", 0)
}

func (suite *Pmon3CoreTestSuite) TestI_StartingAndStopping() {
	suite.cliHelper.ExecBase0("drop")
	suite.cliHelper.ExecCmd("/test/app/bin/test_app", "{\"name\": \"test-server-5\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "running", 0)
	suite.cliHelper.ExecBase1("stop", "1")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "stopped", 0)
	suite.cliHelper.ExecBase2("restart", "1", "{}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(1, "running", 0)
	suite.cliHelper.ExecBase0("drop")
}

func (suite *Pmon3CoreTestSuite) TearDownSuite() {
	god.Banish()
	base.CloseSender()
}
