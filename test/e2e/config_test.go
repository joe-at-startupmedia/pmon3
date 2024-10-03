package e2e

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli"
	"pmon3/conf"
	"pmon3/test/e2e/cli_helper"
	"testing"
	"time"
)

// Define the suite, and absorb the built-in core suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3ConfigTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3ConfigTestSuite))
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Pmon3ConfigTestSuite) SetupSuite() {

	projectPath := os.Getenv("PROJECT_PATH")
	suite.cliHelper = cli_helper.New(&suite.Suite, projectPath)
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3ConfigTestSuite) TestA1_TestConfigFileGetter() {

	assert.Equal(suite.T(), "/etc/pmon3/config/config.yml", conf.GetConfigFile())

	os.Setenv("PMON3_CONF", "/tmp/custom-config.yml")

	assert.Equal(suite.T(), "/tmp/custom-config.yml", conf.GetConfigFile())
}

func (suite *Pmon3ConfigTestSuite) TestA2_TestConfigFileGetter() {

	assert.Equal(suite.T(), "", conf.GetProcessConfigFile())

	os.Setenv("PMON3_PROCESS_CONF", "/tmp/process.custom-config.json")

	assert.Equal(suite.T(), "", conf.GetProcessConfigFile())

	os.Setenv("PMON3_PROCESS_CONF", suite.cliHelper.ProjectPath+"/test/e2e/config/process.core-test.config.json")

	assert.Equal(suite.T(), suite.cliHelper.ProjectPath+"/test/e2e/config/process.core-test.config.json", conf.GetProcessConfigFile())
}

func (suite *Pmon3ConfigTestSuite) TestB_BootCliWithNonExistentConfigFile() {

	configFile := suite.cliHelper.ProjectPath + "/test/e2e/config/nonexistent-test-config.yml"

	if err := cli.Instance(configFile); err != nil {
		suite.FailNow(err.Error())
	}

	assert.Equal(suite.T(), logrus.InfoLevel, cli.Config.GetLogLevel())
}

func (suite *Pmon3ConfigTestSuite) TestC_BootCliWithTestConfigFile() {

	configFile := suite.cliHelper.ProjectPath + "/test/e2e/config/test-config.config.yml"

	if err := cli.Instance(configFile); err != nil {
		suite.FailNow(err.Error())
	}

	assert.Equal(suite.T(), logrus.WarnLevel, cli.Config.GetLogLevel())

	assert.Equal(suite.T(), 1500*time.Millisecond, cli.Config.GetCmdExecResponseWait())

	assert.Equal(suite.T(), 200*time.Millisecond, cli.Config.GetIpcConnectionWait())

	assert.Equal(suite.T(), 1*time.Second, cli.Config.GetDependentProcessEnqueuedWait())

	assert.Equal(suite.T(), 30*time.Second, cli.Config.GetInitializationPeriod())

	os.Setenv("PMON3_DEBUG", "true")

	assert.Equal(suite.T(), logrus.DebugLevel, cli.Config.GetLogLevel())

	os.Setenv("PMON3_DEBUG", "error")

	assert.Equal(suite.T(), logrus.ErrorLevel, cli.Config.GetLogLevel())

	os.Setenv("PMON3_DEBUG", "warn")

	assert.Equal(suite.T(), logrus.WarnLevel, cli.Config.GetLogLevel())

	os.Setenv("PMON3_DEBUG", "info")

	assert.Equal(suite.T(), logrus.InfoLevel, cli.Config.GetLogLevel())

	os.Setenv("PMON3_DEBUG", "debug")

	assert.Equal(suite.T(), logrus.DebugLevel, cli.Config.GetLogLevel())
}

func (suite *Pmon3ConfigTestSuite) TearDownSuite() {}
