package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/test/e2e/cli_helper"
	"testing"
	"time"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3DependencyTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDependencyTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3DependencyTestSuite))
}

func (suite *Pmon3DependencyTestSuite) SetupSuite() {
	suite.cliHelper = cli_helper.SetupSuite(&suite.Suite, "/test/e2e/config/test-config.core.yml", "/test/e2e/config/process.dependency-test.config.json", "dependency")
}

//Alphabetical prefixes are important for ordering: https://github.com/stretchr/testify/issues/194

func (suite *Pmon3DependencyTestSuite) TestA_BootedFromProcessConfigInCorrectOrder() {

	time.Sleep(5 * time.Second)
	passing, cmdResp := suite.cliHelper.LsAssertStatus(5, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	nonDeptProcessNames, deptProcessNames := suite.cliHelper.DgraphProcessNames("")
	assert.Equal(suite.T(), "dep-test-server-5", nonDeptProcessNames[0])
	suite.assertProcessOrder([][]string{
		{"dep-test-server-3"},
		{"dep-test-server-4"},
		{"dep-test-server-2"},
		{"dep-test-server-1"},
	}, deptProcessNames)

	suite.assertProcessOrderFromCmdResp([]string{
		"dep-test-server-3",
		"dep-test-server-4",
		"dep-test-server-2",
		"dep-test-server-1",
		"dep-test-server-5",
	}, processList, "running")
}

func (suite *Pmon3DependencyTestSuite) TestB_AddingAdditionalProcessesWithDeps() {

	execFlags := model.ExecFlags{
		Name:         "dep-test-server-6",
		EnvVars:      "TEST_APP_PORT=11008",
		Dependencies: []string{"dep-test-server-7"},
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", execFlags.Json())

	time.Sleep(2 * time.Second)

	//this should be overwritten by the process configuration file on the next test
	execFlags = model.ExecFlags{
		Dependencies: []string{"dep-test-server-6"},
	}

	suite.cliHelper.ExecBase2("restart", "1", execFlags.Json())

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(6, "running", 0)

	if !passing {
		return
	}

	processList := cmdResp.GetProcessList().GetProcesses()

	//for i := range processList {
	//	pl := processList[i]
	//	log.Printf("%-v", pl)
	//}

	nonDeptProcessNames, deptProcessNames := suite.cliHelper.DgraphProcessNames("")
	assert.Equal(suite.T(), "dep-test-server-5", nonDeptProcessNames[0])
	suite.assertProcessOrder([][]string{
		{"dep-test-server-3"},
		{"dep-test-server-6", "dep-test-server-4"},
		{"dep-test-server-2"},
		{"dep-test-server-1"},
	}, deptProcessNames)

	passing = suite.assertProcessOrderFromCmdResp([]string{
		"dep-test-server-3",
		"dep-test-server-4",
		"dep-test-server-2",
		"dep-test-server-1",
		"dep-test-server-5",
		"dep-test-server-6",
	}, processList, "running")

	if !passing {
		return
	}

	execFlags = model.ExecFlags{
		Name:    "dep-test-server-7",
		EnvVars: "TEST_APP_PORT=11009",
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", execFlags.Json())

	time.Sleep(2 * time.Second)
	suite.cliHelper.LsAssertStatus(7, "running", 0)

	nonDeptProcessNames, deptProcessNames = suite.cliHelper.DgraphProcessNames("")
	assert.Equal(suite.T(), "dep-test-server-5", nonDeptProcessNames[0])
	suite.assertProcessOrder([][]string{
		{"dep-test-server-7", "dep-test-server-3"},
		{"dep-test-server-6", "dep-test-server-4"},
		{"dep-test-server-2"},
		{"dep-test-server-1"},
	}, deptProcessNames)
}

func (suite *Pmon3DependencyTestSuite) TestC_ShouldRebootWithCorrectDependencyOrder() {

	passing := suite.cliHelper.ShouldKill(7, 3)

	if !passing {
		return
	}

	suite.cliHelper.ExecBase2("init", "", "blocking")
	time.Sleep(3 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}

	processList := cmdResp.GetProcessList().GetProcesses()

	//for i := range processList {
	//	pl := processList[i]
	//	log.Printf("%-v", pl)
	//}

	nonDeptProcessNames, deptProcessNames := suite.cliHelper.DgraphProcessNames("")
	assert.Equal(suite.T(), "dep-test-server-5", nonDeptProcessNames[0])
	suite.assertProcessOrder([][]string{
		{"dep-test-server-3", "dep-test-server-7"},
		{"dep-test-server-6", "dep-test-server-4"},
		{"dep-test-server-2"},
		{"dep-test-server-1"},
	}, deptProcessNames)

	passing = suite.assertProcessOrderFromCmdResp([]string{
		"dep-test-server-3",
		"dep-test-server-4",
		"dep-test-server-2",
		"dep-test-server-1",
		"dep-test-server-5",
		"dep-test-server-6",
		"dep-test-server-7",
	}, processList, "running")

}

func (suite *Pmon3DependencyTestSuite) TestD_ShouldRebootFromConfigOnlyWithCorrectDependencyOrder() {

	passing := suite.cliHelper.ShouldKill(7, 3)

	if !passing {
		return
	}

	suite.cliHelper.ExecBase2("init", "process-config-only", "blocking")
	time.Sleep(3 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(5, "running", 0)

	if !passing {
		return
	}

	processList := cmdResp.GetProcessList().GetProcesses()

	//for i := range processList {
	//	pl := processList[i]
	//	log.Printf("%-v", pl)
	//}

	nonDeptProcessNames, deptProcessNames := suite.cliHelper.DgraphProcessNames("process-config-only")
	assert.Equal(suite.T(), "dep-test-server-5", nonDeptProcessNames[0])
	suite.assertProcessOrder([][]string{
		{"dep-test-server-3"},
		{"dep-test-server-4"},
		{"dep-test-server-2"},
		{"dep-test-server-1"},
	}, deptProcessNames)

	suite.assertProcessOrderFromCmdResp([]string{
		"dep-test-server-3",
		"dep-test-server-4",
		"dep-test-server-2",
		"dep-test-server-1",
		"dep-test-server-5",
	}, processList, "running")

}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3DependencyTestSuite) TestZ_TearDown() {
	suite.cliHelper.DropAndClose()
}

func (suite *Pmon3DependencyTestSuite) assertProcessOrderFromCmdResp(processNames []string, processList []*protos.Process, expectedStatus string) bool {

	matchingProcessLen := 0
	for i := range processNames {
		pn := processNames[i]
		if len(expectedStatus) > 0 && processList[i].GetStatus() != expectedStatus {
			continue
		}
		passing := assert.Equal(suite.T(), pn, processList[i].GetName())
		if !passing {
			return false
		}

		matchingProcessLen++
	}

	return assert.Equal(suite.T(), len(processNames), matchingProcessLen)
}

// this is necessary because dgraph is nondeterministic
func (suite *Pmon3DependencyTestSuite) assertProcessOrder(expectedProcessNames [][]string, actualProcessNames []string) bool {
	expectedProcessNameLen := 0
	for i := range expectedProcessNames {
		for range expectedProcessNames[i] {
			expectedProcessNameLen++
		}
	}

	passing := assert.Equal(suite.T(), expectedProcessNameLen, len(actualProcessNames))
	if !passing {
		return false
	}
	assertionsPassed := true
	index := 0
	for k := range expectedProcessNames {
		expectedProcessNameRow := expectedProcessNames[k]
		for range expectedProcessNameRow {
			log.Printf("Does %-v contain %s", expectedProcessNameRow, actualProcessNames[index])
			assert.Contains(suite.T(), expectedProcessNameRow, actualProcessNames[index])
			index++
			if !passing {
				assertionsPassed = false
			}
		}
	}
	return assertionsPassed
}
