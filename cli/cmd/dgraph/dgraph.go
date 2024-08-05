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
	queueOrder := response[0]
	dGraph := response[1]

	fmt.Println("Queue Order")
	for i, appName := range strings.Split(queueOrder, "\n") {
		fmt.Printf("\t%d: %s\n", i, appName)
	}

	if len(dGraph) > 0 {
		fmt.Println("Dependency Graph Order")
		for i, appName := range strings.Split(dGraph, "\n") {
			fmt.Printf("\t%d: %s\n", i, appName)
		}
	}
}
