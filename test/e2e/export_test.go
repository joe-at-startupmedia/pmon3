package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"pmon3/cli/cmd"
	"pmon3/test/e2e/cli_helper"
	"testing"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3ExportTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExportTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3ExportTestSuite))
}

func (suite *Pmon3ExportTestSuite) SetupSuite() {
	suite.cliHelper = cli_helper.SetupSuite(&suite.Suite, "/test/e2e/config/test-config.core.yml", "/test/e2e/config/export/process.from.json", "export")
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3ExportTestSuite) TestA_BootedFromProcessConfig() {
	suite.cliHelper.LsAssertStatus(4, "running", 0)
}

func (suite *Pmon3ExportTestSuite) TestB_ExportJson() {

	exportString, err := cmd.GetExportString("json", "name")
	if err != nil {
		suite.Fail(err.Error())
	}

	fileContents := suite.getExpectedFileContents("process.expected.json")

	assert.Equal(suite.T(), exportString, fileContents)
}

func (suite *Pmon3ExportTestSuite) TestC_ExportToml() {

	exportString, err := cmd.GetExportString("toml", "name")
	if err != nil {
		suite.Fail(err.Error())
	}

	fmt.Println(exportString)

	fileContents := suite.getExpectedFileContents("process.expected.toml")

	assert.Equal(suite.T(), exportString, fileContents)
}

func (suite *Pmon3ExportTestSuite) TestD_ExportYaml() {

	exportString, err := cmd.GetExportString("yaml", "name")
	if err != nil {
		suite.Fail(err.Error())
	}

	fmt.Println(exportString)

	fileContents := suite.getExpectedFileContents("process.expected.yaml")

	assert.Equal(suite.T(), exportString, fileContents)
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3ExportTestSuite) TestZ_TearDown() {
	suite.cliHelper.DropAndClose()
}

func (suite *Pmon3ExportTestSuite) getExpectedFileContents(expectedFile string) string {
	processConfigFile := suite.cliHelper.ProjectPath + "/test/e2e/config/export/" + expectedFile
	fileContents, fileErr := os.ReadFile(processConfigFile)
	if fileErr != nil {
		suite.Fail(fileErr.Error())
	}
	return string(fileContents)
}
