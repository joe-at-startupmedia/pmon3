package e2e

// Basic imports
import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/conf"
	"pmon3/pmond/model"
	"testing"

	"pmon3/pmond/protos"

	"time"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3BasicTestSuite struct {
	suite.Suite
	appBinPath string
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

	suite.appBinPath = os.Getenv("APP_BIN_PATH")

	base.OpenSender()
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3BasicTestSuite) TestA_BootedFromProcessConfig() {
	suite.lsAssert(2)
}

func (suite *Pmon3BasicTestSuite) TestB_AddingAdditionalProcessesFromProcessConfig() {
	suite.execCmd("/bin/app", "{\"name\": \"test-server3\"}")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(3, "running", 0)

	suite.execCmd("/bin/app", "{\"name\": \"test-server4\"}")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(4, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestC_DescribingAProcessWithAFourthId() {
	newCmdResp := suite.execBase1("desc", "4")
	assert.Equal(suite.T(), "test-server4", newCmdResp.GetProcess().GetName())

}

func (suite *Pmon3BasicTestSuite) TestD_DeletingAProcess() {
	suite.execBase1("del", "3")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(3, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestE_KillProcesses() {
	var sent *protos.Cmd
	sent = base.SendCmd("kill", "")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		time.Sleep(2 * time.Second)
		suite.lsAssertStatus(3, "stopped", 0)
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
		suite.lsAssertStatus(3, "running", 0)
	}
}

func (suite *Pmon3BasicTestSuite) TestG_Drop() {
	suite.execBase0("drop")
	time.Sleep(2 * time.Second)
	suite.lsAssert(0)
}

func (suite *Pmon3BasicTestSuite) TestH_InitAllAfterDrop() {
	suite.execBase2("init", "all", "blocking")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(2, "running", 0)
}

func (suite *Pmon3BasicTestSuite) TestI_StartingAndStopping() {
	suite.execBase0("drop")
	suite.execCmd("/bin/app", "{\"name\": \"test-server5\"}")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(1, "running", 0)
	suite.execBase1("stop", "1")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(1, "stopped", 0)
	suite.execBase2("restart", "1", "{}")
	time.Sleep(2 * time.Second)
	suite.lsAssertStatus(1, "running", 0)
	suite.execBase0("drop")
}

func (suite *Pmon3BasicTestSuite) TearDownSuite() {
	base.CloseSender()
}

func (suite *Pmon3BasicTestSuite) lsAssert(expectedProcessLen int) *protos.CmdResp {
	newCmdResp := suite.execBase0("list")
	processList := newCmdResp.GetProcessList().GetProcesses()
	cli.Log.Infof("process list: %s \n value string: %s \n", processList, newCmdResp.GetValueStr())
	assert.Equal(suite.T(), expectedProcessLen, len(processList))
	//cli.Log.Fatalf("Expected process length of %d but got %d", expectedProcessLen, actualProcessLen)
	return newCmdResp
}

func (suite *Pmon3BasicTestSuite) lsAssertStatus(expectedProcessLen int, status string, retries int) {

	cmdResp := suite.lsAssert(expectedProcessLen)
	processList := cmdResp.GetProcessList().GetProcesses()

	for _, p := range processList {
		if p.Status != status && retries < 3 { //three retries are allowed
			cli.Log.Infof("Expected process status of %s but got %s", status, p.Status)
			cli.Log.Warnf("retry count: %d", retries+1)
			time.Sleep(time.Second * 5)
			suite.lsAssertStatus(expectedProcessLen, status, retries+1)
			break
		} else {
			assert.Equal(suite.T(), status, p.Status)
		}
	}
}

func (suite *Pmon3BasicTestSuite) execCmd(processFile string, execFlagsJson string) *protos.CmdResp {
	processFile = suite.appBinPath + processFile
	cli.Log.Infof("Executing: pmon3 exec %s %s", processFile, execFlagsJson)
	ef := model.ExecFlags{}
	execFlags, err := ef.Parse(execFlagsJson)
	if err != nil {
		cli.Log.Fatal(err)
	}
	execFlags.File = processFile
	return suite.execBase1("exec", execFlags.Json())
}

func (suite *Pmon3BasicTestSuite) execBase0(cmd string) *protos.CmdResp {
	return suite.execBase(cmd, "", "")
}

func (suite *Pmon3BasicTestSuite) execBase1(cmd string, arg1 string) *protos.CmdResp {
	return suite.execBase(cmd, arg1, "")
}

func (suite *Pmon3BasicTestSuite) execBase2(cmd string, arg1 string, arg2 string) *protos.CmdResp {
	return suite.execBase(cmd, arg1, arg2)
}

func (suite *Pmon3BasicTestSuite) execBase(cmd string, arg1 string, arg2 string) *protos.CmdResp {
	var sent *protos.Cmd
	if len(arg2) > 0 {
		sent = base.SendCmdArg2(cmd, arg1, arg2)
	} else {
		sent = base.SendCmd(cmd, arg1)
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	}
	return newCmdResp
}
