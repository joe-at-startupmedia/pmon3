package e2e

// Basic imports
import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
	"pmon3/test/e2e/cli_helper"
	"testing"

	"pmon3/pmond/protos"

	"time"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3BasicTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3BasicTestSuite))
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Pmon3BasicTestSuite) SetupSuite() {
	if err := cli.Instance(conf.GetConfigFile()); err != nil {
		suite.FailNow(err.Error())
	}

	appBinPath := os.Getenv("APP_BIN_PATH")
	suite.cliHelper = cli_helper.New(&suite.Suite, appBinPath)

	base.OpenSender()
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3BasicTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssert(2)
}

func (suite *Pmon3BasicTestSuite) TestB_AddingAdditionalProcessesFromProcessConfig() {
	suite.cliHelper.ExecCmd("/bin/app", "{\"name\": \"test-server3\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(3, "running", 0)

	suite.cliHelper.ExecCmd("/bin/app", "{\"name\": \"test-server4\"}")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(4, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestC_DescribingAProcessWithAFourthId() {
	newCmdResp := suite.cliHelper.ExecBase1("desc", "4")
	assert.Equal(suite.T(), "test-server4", newCmdResp.GetProcess().GetName())

}

func (suite *Pmon3BasicTestSuite) TestD_DeletingAProcess() {
	suite.cliHelper.ExecBase1("del", "3")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(3, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestE_KillProcesses() {
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

func (suite *Pmon3BasicTestSuite) TestF_InitAll() {
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

func (suite *Pmon3BasicTestSuite) TestG_Drop() {
	suite.cliHelper.ExecBase0("drop")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssert(0)
}

func (suite *Pmon3BasicTestSuite) TestH_InitAllAfterDrop() {
	suite.cliHelper.ExecBase2("init", "all", "blocking")
	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(2, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestI_StartingAndStopping() {
	suite.cliHelper.ExecBase0("drop")
	suite.cliHelper.ExecCmd("/bin/app", "{\"name\": \"test-server5\"}")
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

func (suite *Pmon3BasicTestSuite) TearDownSuite() {
	base.CloseSender()
}
