package dgraph

import (
	"fmt"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/protos"
	"strings"

	"github.com/spf13/cobra"
)

var (
	processConfigOnlyFlag bool
)

var Cmd = &cobra.Command{
	Use:     "dgraph",
	Aliases: []string{"order"},
	Short:   "Show the process queue order",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Dgraph(processConfigOnlyFlag)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&processConfigOnlyFlag, "process-config-only", "c", false, "only initialize processes specified in the Processes Config file")
}

func Dgraph(processConfigOnly bool) *protos.CmdResp {

	var sent *protos.Cmd

	if processConfigOnly {
		sent = base.SendCmd("dgraph", "process-config-only")
	} else {
		sent = base.SendCmd("dgraph", "")
	}

	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {

		response := strings.Split(newCmdResp.GetValueStr(), "||")

		var nonDeptProcessNames []string
		var deptProcessNames []string
		if len(response[0]) > 0 {
			nonDeptProcessNames = strings.Split(response[0], "\n")
		}
		if len(response[1]) > 0 {
			deptProcessNames = strings.Split(response[1], "\n")
		}

		fmt.Println("Queue Order")
		for i, processName := range deptProcessNames {
			fmt.Printf("\t%d: %s\n", i, processName)
		}
		for i, processName := range nonDeptProcessNames {
			fmt.Printf("\t%d: %s\n", i+len(deptProcessNames), processName)
		}

		if len(deptProcessNames) > 0 {
			fmt.Println("Dependency Graph Order")
			for i, processName := range deptProcessNames {
				fmt.Printf("\t%d: %s\n", i, processName)
			}
		}
	}

	return newCmdResp
}
