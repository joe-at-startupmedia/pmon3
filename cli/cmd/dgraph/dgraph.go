package dgraph

import (
	"fmt"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"strings"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "dgraph",
	Aliases: []string{"order"},
	Short:   "Show the process queue order",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		Dgraph()
		base.CloseSender()
	},
}

func Dgraph() {

	sent := base.SendCmd("dgraph", "")
	newCmdResp := base.GetResponse(sent)

	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}

	response := strings.Split(newCmdResp.GetValueStr(), "||")

	var nonDeptAppNames []string
	var deptAppNames []string
	if len(response[0]) > 0 {
		nonDeptAppNames = strings.Split(response[0], "\n")
	}
	if len(response[1]) > 0 {
		deptAppNames = strings.Split(response[1], "\n")
	}

	fmt.Println("Queue Order")
	for i, appName := range deptAppNames {
		fmt.Printf("\t%d: %s\n", i, appName)
	}
	for i, appName := range nonDeptAppNames {
		fmt.Printf("\t%d: %s\n", i+len(deptAppNames), appName)
	}

	if len(deptAppNames) > 0 {
		fmt.Println("Dependency Graph Order")
		for i, appName := range deptAppNames {
			fmt.Printf("\t%d: %s\n", i, appName)
		}
	}
}
