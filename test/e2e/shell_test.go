package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"pmon3/cli/shell"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	shell2 "pmon3/pmond/shell"
	"pmon3/test/e2e/cli_helper"
	"testing"
	"time"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3ShellTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestShellTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3ShellTestSuite))
}

func (suite *Pmon3ShellTestSuite) SetupSuite() {
	suite.cliHelper = cli_helper.SetupSuite(&suite.Suite, "/test/e2e/config/test-config.core.yml", "", "shell")
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3ShellTestSuite) TestA_PmondIsntRunningFromCli() {
	isRunning := shell.ExecIsPmondRunning()
	assert.Equal(suite.T(), isRunning, false)
}

func (suite *Pmon3ShellTestSuite) TestB_PmondIsntRunningFromPmond() {
	isRunning := shell2.ExecIsPmondRunning(7777) //use a pid number unlikely to exist
	assert.Equal(suite.T(), isRunning, false)
}

var pid int

func (suite *Pmon3ShellTestSuite) TestC_PmondIsRunningAfterStarting() {
	cmd := exec.Command("pmond")
	cmd.Env = []string{
		fmt.Sprintf("PMON3_CONF=%s/test/e2e/config/test-config.core.yml", suite.cliHelper.ProjectPath),
	}
	go func() {
		if err := cmd.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	timer := time.NewTimer(time.Second)
	for {
		<-timer.C
		if cmd.Process.Pid > 0 {
			pid = cmd.Process.Pid
			fmt.Println("Got pid", pid)
			return
		}
	}
}

func (suite *Pmon3ShellTestSuite) TestD_PmondIsRunningFromCli() {
	isRunning := shell.ExecIsPmondRunning()
	assert.Equal(suite.T(), isRunning, true)
}

func (suite *Pmon3ShellTestSuite) TestE_PmondIsRunningFromPmond() {
	isRunning := shell2.ExecIsPmondRunning(os.Getpid()) //use a pid number unlikely to exist
	assert.Equal(suite.T(), isRunning, true)
}

func (suite *Pmon3ShellTestSuite) TestF_KillPid() {
	if err := process.SendOsKillSignal(&model.Process{
		Pid: uint32(pid),
	}, true); err != nil {
		suite.Fail(err.Error())
	}
}

func (suite *Pmon3ShellTestSuite) TestG_PmondIsntRunningFromCli() {
	time.Sleep(time.Second * 5)
	isRunning := shell.ExecIsPmondRunning()
	assert.Equal(suite.T(), isRunning, false)
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3ShellTestSuite) TestZ_TearDown() {
	suite.cliHelper.Close()
}
