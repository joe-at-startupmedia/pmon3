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
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/array"
	"pmon3/test/e2e/cli_helper"
	"testing"
	"time"
)

// Define the suite, and absorb the built-in suite
// functionality from testify - including a T() method which
// returns the current testing context
type Pmon3GroupTestSuite struct {
	suite.Suite
	cliHelper *cli_helper.CliHelper
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGroupTestSuite(t *testing.T) {
	suite.Run(t, new(Pmon3GroupTestSuite))
}

func (suite *Pmon3GroupTestSuite) SetupSuite() {
	projectPath := os.Getenv("PROJECT_PATH")
	suite.cliHelper = cli_helper.New(&suite.Suite, projectPath)

	configFile := projectPath + "/test/e2e/config/test-config.core.yml"
	processConfigFile := projectPath + "/test/e2e/config/process.group-test.config.json"
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

func (suite *Pmon3GroupTestSuite) TestA1_BootedFromProcessConfigWithCorrectGroups() {

	time.Sleep(5 * time.Second)
	passing, cmdResp := suite.cliHelper.LsAssertStatus(5, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	var groupNames = array.Map(processList, func(p *protos.Process) []string {
		return array.Map(p.GetGroups(), func(g *protos.Group) string {
			return g.GetName()
		})
	})

	assert.Equal(suite.T(), len(groupNames), 5)
	assert.Equal(suite.T(), []string{"groupA"}, groupNames[0])
	assert.Equal(suite.T(), []string{"groupB"}, groupNames[1])
	assert.Equal(suite.T(), []string{"groupA", "groupB"}, groupNames[2])
	assert.Equal(suite.T(), []string{"groupC"}, groupNames[3])
	assert.Empty(suite.T(), groupNames[4])

}

func (suite *Pmon3GroupTestSuite) TestA2_ListGroups() {

	cmdResp := suite.cliHelper.ExecBase0("group_list")

	groupList := cmdResp.GetGroupList().GetGroups()

	assert.Equal(suite.T(), 3, len(groupList))
	assert.Equal(suite.T(), "groupA", groupList[0].GetName())
	assert.Equal(suite.T(), "groupB", groupList[1].GetName())
	assert.Equal(suite.T(), "groupC", groupList[2].GetName())
}

func (suite *Pmon3GroupTestSuite) TestA3_GetProcessGroups() {

	cmdResp := suite.cliHelper.ExecBase1("desc", "group-test-server-3")

	groupList := cmdResp.GetProcess().GetGroups()

	assert.Equal(suite.T(), 2, len(groupList))
	assert.Equal(suite.T(), "groupA", groupList[0].GetName())
	assert.Equal(suite.T(), "groupB", groupList[1].GetName())
}

func (suite *Pmon3GroupTestSuite) TestB1_ExecCmdWithNewGroup() {
	execFlags := model.ExecFlags{
		Name:    "group-test-server-6",
		EnvVars: "TEST_APP_PORT=11015",
		Groups:  []string{"groupD"},
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", execFlags.Json())

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(6, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), []string{"groupD"}, array.Map(processList[5].GetGroups(), func(g *protos.Group) string {
		return g.GetName()
	}))
}

func (suite *Pmon3GroupTestSuite) TestB2_ListGroups() {

	cmdResp := suite.cliHelper.ExecBase0("group_list")

	groupList := cmdResp.GetGroupList().GetGroups()

	assert.Equal(suite.T(), 4, len(groupList))
	assert.Equal(suite.T(), "groupA", groupList[0].GetName())
	assert.Equal(suite.T(), "groupB", groupList[1].GetName())
	assert.Equal(suite.T(), "groupC", groupList[2].GetName())
	assert.Equal(suite.T(), "groupD", groupList[3].GetName())
}

func (suite *Pmon3GroupTestSuite) TestC1_ExecCmdWithExistingGroup() {
	execFlags := model.ExecFlags{
		Name:    "group-test-server-7",
		EnvVars: "TEST_APP_PORT=11016",
		Groups:  []string{"groupC", "groupE"},
	}

	suite.cliHelper.ExecCmd("/test/app/bin/test_app", execFlags.Json())

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), []string{"groupC", "groupE"}, array.Map(processList[6].GetGroups(), func(g *protos.Group) string {
		return g.GetName()
	}))
}

func (suite *Pmon3GroupTestSuite) TestC2_ListGroups() {

	cmdResp := suite.cliHelper.ExecBase0("group_list")

	groupList := cmdResp.GetGroupList().GetGroups()

	assert.Equal(suite.T(), 5, len(groupList))

	assert.Equal(suite.T(), "groupA", groupList[0].GetName())
	assert.Equal(suite.T(), "groupB", groupList[1].GetName())
	assert.Equal(suite.T(), "groupC", groupList[2].GetName())
	assert.Equal(suite.T(), "groupD", groupList[3].GetName())
	assert.Equal(suite.T(), "groupE", groupList[4].GetName())
}

func (suite *Pmon3GroupTestSuite) TestD_RestartGroupA() {

	suite.cliHelper.ExecBase2("group_restart", "groupA", "{}")

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), uint32(1), processList[0].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[1].GetRestartCount())
	assert.Equal(suite.T(), uint32(1), processList[2].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[3].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[4].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[5].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[6].GetRestartCount())

}

func (suite *Pmon3GroupTestSuite) TestE_RestartGroupB() {

	suite.cliHelper.ExecBase2("group_restart", "groupB", "{}")

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), uint32(1), processList[0].GetRestartCount())
	assert.Equal(suite.T(), uint32(1), processList[1].GetRestartCount())
	assert.Equal(suite.T(), uint32(2), processList[2].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[3].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[4].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[5].GetRestartCount())
	assert.Equal(suite.T(), uint32(0), processList[6].GetRestartCount())
}

func (suite *Pmon3GroupTestSuite) TestE_StopGroupA() {

	suite.cliHelper.ExecBase1("group_stop", "groupA")

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(5, "running", 0)

	if !passing {
		return
	}

	passing, cmdResp = suite.cliHelper.LsAssertStatus(2, "stopped", 0)

	if !passing {
		return
	}

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), model.StatusStopped.String(), processList[0].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[1].GetStatus())
	assert.Equal(suite.T(), model.StatusStopped.String(), processList[2].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[3].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[4].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[5].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[6].GetStatus())

}

func (suite *Pmon3GroupTestSuite) TestE_StopGroupC() {

	suite.cliHelper.ExecBase1("group_stop", "groupC")

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(3, "running", 0)

	if !passing {
		return
	}

	passing, cmdResp = suite.cliHelper.LsAssertStatus(4, "stopped", 0)

	if !passing {
		return
	}

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), model.StatusStopped.String(), processList[0].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[1].GetStatus())
	assert.Equal(suite.T(), model.StatusStopped.String(), processList[2].GetStatus())
	assert.Equal(suite.T(), model.StatusStopped.String(), processList[3].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[4].GetStatus())
	assert.Equal(suite.T(), model.StatusRunning.String(), processList[5].GetStatus())
	assert.Equal(suite.T(), model.StatusStopped.String(), processList[6].GetStatus())

}

func (suite *Pmon3GroupTestSuite) TestF_CreateGroup() {

	suite.cliHelper.ExecBase1("group_create", "groupF")

	time.Sleep(2 * time.Second)

	cmdResp := suite.cliHelper.ExecBase0("group_list")

	groupList := cmdResp.GetGroupList().GetGroups()

	assert.Equal(suite.T(), 6, len(groupList))

	assert.Equal(suite.T(), "groupA", groupList[0].GetName())
	assert.Equal(suite.T(), "groupB", groupList[1].GetName())
	assert.Equal(suite.T(), "groupC", groupList[2].GetName())
	assert.Equal(suite.T(), "groupD", groupList[3].GetName())
	assert.Equal(suite.T(), "groupE", groupList[4].GetName())
	assert.Equal(suite.T(), "groupF", groupList[5].GetName())
}

func (suite *Pmon3GroupTestSuite) TestG_AssignGroup() {

	suite.cliHelper.ExecBase2("group_assign", "groupF", "group-test-server-5")

	time.Sleep(2 * time.Second)

	cmdResp := suite.cliHelper.ExecBase1("group_desc", "groupF")

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), len(processList), 1)

	assert.Equal(suite.T(), "group-test-server-5", processList[0].GetName())
}

func (suite *Pmon3GroupTestSuite) TestH_RemoveGroup() {

	suite.cliHelper.ExecBase2("group_remove", "groupF", "group-test-server-5")

	time.Sleep(2 * time.Second)

	cmdResp := suite.cliHelper.ExecBase1("group_desc", "groupF")

	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), len(processList), 0)
}

func (suite *Pmon3GroupTestSuite) TestI1_DeleteGroup() {

	suite.cliHelper.ExecBase2("group_restart", "groupA", "{}")

	time.Sleep(2 * time.Second)

	suite.cliHelper.ExecBase2("group_restart", "groupC", "{}")

	time.Sleep(2 * time.Second)

	passing, cmdResp := suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}

	suite.cliHelper.ExecBase1("group_del", "groupA")

	time.Sleep(2 * time.Second)

	passing, cmdResp = suite.cliHelper.LsAssertStatus(7, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	var groupNames = array.Map(processList, func(p *protos.Process) []string {
		return array.Map(p.GetGroups(), func(g *protos.Group) string {
			return g.GetName()
		})
	})

	//deleting a group doesn't delete the process
	assert.Equal(suite.T(), 7, len(groupNames))
	assert.Empty(suite.T(), groupNames[0])
	assert.Equal(suite.T(), []string{"groupB"}, groupNames[1])
	assert.Equal(suite.T(), []string{"groupB"}, groupNames[2])
	assert.Equal(suite.T(), []string{"groupC"}, groupNames[3])
	assert.Empty(suite.T(), groupNames[4])
	assert.Equal(suite.T(), []string{"groupD"}, groupNames[5])
	assert.Equal(suite.T(), []string{"groupC", "groupE"}, groupNames[6])
}

func (suite *Pmon3GroupTestSuite) TestJ1_DropGroup() {

	cmdResp := suite.cliHelper.ExecBase1("group_drop", "groupB")

	passing, cmdResp := suite.cliHelper.LsAssertStatus(5, "running", 0)

	if !passing {
		return
	}
	processList := cmdResp.GetProcessList().GetProcesses()

	assert.Equal(suite.T(), 5, len(processList))

	var groupNames = array.Map(processList, func(p *protos.Process) []string {
		return array.Map(p.GetGroups(), func(g *protos.Group) string {
			return g.GetName()
		})
	})

	//deleting a group doesn't delete the process
	assert.Equal(suite.T(), 5, len(groupNames))
	assert.Empty(suite.T(), groupNames[0])
	assert.Equal(suite.T(), []string{"groupC"}, groupNames[1])
	assert.Empty(suite.T(), groupNames[2])
	assert.Equal(suite.T(), []string{"groupD"}, groupNames[3])
	assert.Equal(suite.T(), []string{"groupC", "groupE"}, groupNames[4])
}

func (suite *Pmon3GroupTestSuite) TestJ2_ListGroups() {

	cmdResp := suite.cliHelper.ExecBase0("group_list")

	groupList := cmdResp.GetGroupList().GetGroups()

	assert.Equal(suite.T(), 5, len(groupList))

	//dropping a group doesnt delete it
	assert.Equal(suite.T(), "groupB", groupList[0].GetName())
	assert.Equal(suite.T(), "groupC", groupList[1].GetName())
	assert.Equal(suite.T(), "groupD", groupList[2].GetName())
	assert.Equal(suite.T(), "groupE", groupList[3].GetName())
	assert.Equal(suite.T(), "groupF", groupList[4].GetName())
}

// this is necessary because TearDownSuite executes concurrently with the
// initialization of the next suite
func (suite *Pmon3GroupTestSuite) TestZ_TearDown() {
	suite.cliHelper.DropAndClose()
}
