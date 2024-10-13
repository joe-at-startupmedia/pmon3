package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	initialize "pmon3/cli/cmd"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/flap_detector"
	"pmon3/pmond/observer"
	"pmon3/test/e2e/cli_helper"
	"testing"
	"time"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3FlapTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

var eventCounters = map[string]int{}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFlapTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3FlapTestSuite))
}

func (suite *Pmon3FlapTestSuite) SetupSuite() {
	suite.cliHelper = cli_helper.SetupSuite(&suite.Suite, "/test/e2e/config/test-config.flap.yml", "/test/e2e/config/process.flap-test.config.json", "flap")
	observer.OnRestartEventFunc = func(evt *observer.Event) {
		fmt.Println("OnRestartEventFunc incrementing")
		eventCounters["restart"]++
	}
	observer.OnFailedEventFunc = func(evt *observer.Event) {
		fmt.Println("OnFailedEventFunc incrementing")
		eventCounters["failed"]++
	}
	observer.OnBackOffEventFunc = func(evt *observer.Event) {
		fmt.Println("OnBackOffEventFunc incrementing")
		eventCounters["backoff"]++
	}
	//prevent the controller from reloading the config file
	pmond.Config.DisableReloads = true
	pmond.Config.EventHandler = conf.EventHandlerConfig{
		ProcessRestart: suite.cliHelper.ProjectPath + "/test/e2e/config/on_event.bash",
		ProcessFailure: suite.cliHelper.ProjectPath + "/test/e2e/config/on_event.bash",
		ProcessBackoff: suite.cliHelper.ProjectPath + "/test/e2e/config/on_event.bash",
	}
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3FlapTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssert(4)
}

func (suite *Pmon3FlapTestSuite) TestB_ShouldBackoff() {

	passing, _ := suite.cliHelper.LsAssertStatus(2, "running", 0)
	if !passing {
		return
	}

	for range 5 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be failed with %d restarts\n", 2)
	newCmdResp := suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), uint32(2), newCmdResp.GetProcess().GetRestartCount())
	assert.Equal(suite.T(), "failed", newCmdResp.GetProcess().GetStatus())

	for range 3 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 3)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), uint32(3), newCmdResp.GetProcess().GetRestartCount())
	assert.Equal(suite.T(), "backoff", newCmdResp.GetProcess().GetStatus())
	assert.Equal(suite.T(), 1, eventCounters["backoff"])
}

func (suite *Pmon3FlapTestSuite) TestC_ShouldBackoffWith3Restarts() {

	for range 8 {
		initialize.List()
		time.Sleep(time.Second * 2)
		fmt.Println(flap_detector.FromProcessId(2, pmond.Config))
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 4)
	newCmdResp := suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), uint32(4), newCmdResp.GetProcess().GetRestartCount())
	assert.Equal(suite.T(), "backoff", newCmdResp.GetProcess().GetStatus())
	assert.Equal(suite.T(), 2, eventCounters["backoff"])
	assert.Equal(suite.T(), 1, eventCounters["restart"])

	for range 8 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 5)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), uint32(5), newCmdResp.GetProcess().GetRestartCount())
	assert.Equal(suite.T(), "backoff", newCmdResp.GetProcess().GetStatus())
	assert.Equal(suite.T(), 3, eventCounters["backoff"])
	assert.Equal(suite.T(), 2, eventCounters["restart"])

	for range 8 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 6)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(6))
	assert.Equal(suite.T(), "backoff", newCmdResp.GetProcess().GetStatus())
	assert.Equal(suite.T(), 4, eventCounters["backoff"])
	assert.Equal(suite.T(), 3, eventCounters["restart"])

	suite.cliHelper.ExecBase1("drop", "force")
	time.Sleep(5 * time.Second)
}

func (suite *Pmon3FlapTestSuite) TestD_ShouldBackoff() {

	eventCounters = map[string]int{}
	flap_detector.Reset()
	pmond.Config.FlapDetection.ThresholdDecrement = 0
	pmond.Config.FlapDetection.ThresholdCountdown = 60
	fmt.Println(pmond.Config.Yaml())
	newCmdResp := initialize.Initialize(false, true)
	if len(newCmdResp.GetError()) > 0 {
		suite.Fail(newCmdResp.GetError())
	} else {
		time.Sleep(2 * time.Second)
		suite.cliHelper.LsAssertStatus(2, "running", 0)
	}

	passing, _ := suite.cliHelper.LsAssertStatus(2, "running", 0)
	if !passing {
		return
	}

	for range 5 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be failed with %d restarts\n", 2)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(2))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "failed")
	assert.Equal(suite.T(), 0, eventCounters["backoff"])
	assert.Equal(suite.T(), 2, eventCounters["restart"])

	for range 3 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 3)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(3))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "backoff")
	assert.Equal(suite.T(), 1, eventCounters["backoff"])
	assert.Equal(suite.T(), 3, eventCounters["restart"])
}

func (suite *Pmon3FlapTestSuite) TestE_ShouldBackoffWith3Restarts() {

	for range 8 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 3)
	newCmdResp := suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(3))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "backoff")
	assert.Equal(suite.T(), 1, eventCounters["backoff"])
	assert.Equal(suite.T(), 3, eventCounters["restart"])

	for range 7 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be failed with %d restarts\n", 4)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(4))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "failed")
	assert.Equal(suite.T(), 1, eventCounters["backoff"])
	assert.Equal(suite.T(), 4, eventCounters["restart"])

	for range 2 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be failed with %d restarts\n", 5)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(5))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "failed")
	assert.Equal(suite.T(), 1, eventCounters["backoff"])
	assert.Equal(suite.T(), 5, eventCounters["restart"])

	for range 8 {
		initialize.List()
		time.Sleep(time.Second * 2)
	}

	fmt.Printf("Asserting that process should be backed off with %d restarts\n", 6)
	newCmdResp = suite.cliHelper.ExecBase1("desc", "test-server-21")
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetRestartCount(), uint32(6))
	assert.Equal(suite.T(), newCmdResp.GetProcess().GetStatus(), "backoff")
	assert.Equal(suite.T(), 2, eventCounters["backoff"])
	assert.Equal(suite.T(), 6, eventCounters["restart"])
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3FlapTestSuite) TestZ_TearDown() {
	fmt.Printf("event counters: %-v", eventCounters)
	suite.cliHelper.DropAndClose()
}
